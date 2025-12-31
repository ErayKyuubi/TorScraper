package main

import (
	"TorScraper/internal/input_handler"
	"TorScraper/internal/logger"
	"TorScraper/internal/proxy"
	"TorScraper/internal/scanner"
	"fmt"
	"log"
)

func main() {
	// YAML dosyasını oku
	targets, err := input_handler.ReadTargets("targets.yaml")
	if err != nil {
		logger.LogScanStatus("FILE_READ", "FAILED", err)
		return
	}

	// IP kontrol
	client, _ := proxy.CreateTorClient()

	fmt.Println("[WAIT] Tor IP kontrolü yapılıyor...")
	success, message := proxy.CheckTorConnection(client)

	if !success {
		log.Fatalf("[CRITICAL] %s", message) // Programı durdur
	}
	fmt.Printf("[OK] %s\n", message)

	// Tarama başlat
	scanner.ScanTargets(targets)

	// Bitiş Raporu
	logger.LogScanStatus("ALL TASKS", "SCAN COMPLETE", nil)
	fmt.Println("\n[INFO] HTML dosyaları: output/html/")
	fmt.Println("[INFO] Screenshot'lar: output/screenshots/")
}
