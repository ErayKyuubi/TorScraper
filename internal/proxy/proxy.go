package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"golang.org/x/net/proxy"
)

func CreateTorClient() (*http.Client, error) {
	// Tor Expert Başlat
	cmd := exec.Command("C:\\TorExpert\\tor\\tor.exe")
	err := cmd.Start()
	if err != nil {
		log.Fatal("Tor servisi başlatılamadı:", err)
	}

	// SOCKS5 proxy adresi
	proxyAddr := "127.0.0.1:9050"

	// Dialer
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	// IP sızıntısını önlemek için
	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	// Bu transportu kullanan istemciyi döndür
	return &http.Client{
		Transport: transport,
	}, nil

}

func CheckTorConnection(client *http.Client) (bool, string) {
	// Tor kontrol sayfasına istek at
	resp, err := client.Get("https://check.torproject.org/")
	if err != nil {
		return false, "Bağlantı Hatası: Tor servisi çalışıyor mu?"
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	content := string(body)

	var currentIP string
	// İçerikte IP adresini içeren genel bloğu bul
	if strings.Contains(content, "Your IP address appears to be:") {
		parts := strings.Split(content, "Your IP address appears to be:")
		if len(parts) > 1 {

			temp := strings.Split(parts[1], ">")
			if len(temp) > 1 {
				ipParts := strings.Split(temp[1], "<")
				currentIP = ipParts[0]
			}
		}
	}

	// Sayfa içeriğinde başarılı bağlantı mesajını ara
	if strings.Contains(content, "Congratulations") {
		return true, fmt.Sprintf("Tor bağlantısı başarılı! IP Adresiniz: %s", currentIP)
	}

	return false, fmt.Sprintf("DİKKAT: Tor kullanılmıyor! Görünen IP: %s", currentIP)
}
