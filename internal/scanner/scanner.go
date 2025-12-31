package scanner

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"TorScraper/internal/logger"

	"github.com/chromedp/chromedp"
)

func ScanTargets(urls []string) {
	fmt.Println("[INFO] Tarama işlemi başlatılıyor...")

	numWorkers := 10 // Aynı anda kaç sayfa taranacağını buradan ayarlayın
	jobs := make(chan string, len(urls))
	var wg sync.WaitGroup

	// İşçileri başlat
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for url := range jobs {
				processURL(url)
			}
		}(i)
	}

	// İşleri kanala gönder
	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	// Tüm işçiler bitti
	wg.Wait()
	fmt.Println("[INFO] Tüm tarama işlemleri tamamlandı.")
}

// URL temizleme
func processURL(url string) {
	target := strings.TrimSpace(url)
	if target == "" {
		return
	}

	if !strings.HasPrefix(target, "http") {
		target = "http://" + target
	}

	fmt.Printf("[INFO] Scanning: %s -> IN PROGRESS\n", target)
	err := performScan(target)

	if err != nil {
		logger.LogScanStatus(target, "FAILED", err)
	} else {
		logger.LogScanStatus(target, "SUCCESS", nil)
	}
}

func performScan(targetURL string) error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer("socks5://127.0.0.1:9050"),
		chromedp.NoSandbox,
		chromedp.Headless,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	var screenshotBuf []byte
	var htmlContent string

	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.Sleep(5*time.Second),
		chromedp.FullScreenshot(&screenshotBuf, 90),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		return err
	}

	return saveResults(targetURL, screenshotBuf, htmlContent)
}

func saveResults(url string, screenshot []byte, html string) error {
	os.MkdirAll("output/screenshots", 0755)
	os.MkdirAll("output/html", 0755)

	safeName := strings.ReplaceAll(url, "http://", "")
	safeName = strings.ReplaceAll(safeName, "https://", "")
	safeName = strings.ReplaceAll(safeName, ".", "_")
	safeName = strings.ReplaceAll(safeName, "/", "")

	timestamp := time.Now().Format("20060102_150405_000000") // Çakışma olmasın
	finalBaseName := fmt.Sprintf("%s_%s", safeName, timestamp)

	err := os.WriteFile(fmt.Sprintf("output/screenshots/%s.png", finalBaseName), screenshot, 0644)
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("output/html/%s.html", finalBaseName), []byte(html), 0644)
}
