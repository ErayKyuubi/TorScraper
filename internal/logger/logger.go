package logger

import (
	"fmt"
	"log"
	"os"
)

func LogScanStatus(url string, status string, err error) {
	var message string
	if err != nil {
		message = fmt.Sprintf("[ERR] Scanning: %s -> %v", url, err)
	} else {
		message = fmt.Sprintf("[INFO] Scanning: %s -> %s", url, status)
	}

	fmt.Println(message)

	// scan_report.log dosyasına ekle
	f, _ := os.OpenFile("scan_report.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	log.SetOutput(f)
	log.Println(message)
}
