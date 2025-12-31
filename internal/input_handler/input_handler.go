package input_handler

import (
	"os"
	"strings"
)

func ReadTargets(filePath string) ([]string, error) {
	// Dosyayı oku
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Satırlara böl
	lines := strings.Split(string(data), "\n")

	var targets []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" { // Boş satırları atla
			targets = append(targets, line)
		}
	}

	return targets, nil
}
