package cloudflare

import (
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func Test() {
	getCurrentIP()
}

func getCurrentIP() {
	resp, err := http.Get("https://api.ipify.org")

	if err != nil {
		log.Infof("Error getting IP: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Infof("Failed to read body: %v", err)
	}

	log.Infof("IP: %v", string(body))
}