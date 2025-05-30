package main

import (
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/maxmorhardt/cloudflare-k8s-manager/internal/cloudflare"
	// "github.com/maxmorhardt/cloudflare-k8s-manager/internal/k8s"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)
  	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{                                             
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {                                                     
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)       
			return "", fileName                                                      
		},                                                                           
	})

	err := godotenv.Load()
	if err != nil {
		log.Error("Could not load .env file")
		os.Exit(1)
	}

	config := loadConfig()

	// log.Info("Starting watcher")
	// go k8s.Watcher()
	// go cloudflare.CheckDNSExists(config, "api.maxstash.io")
	go cloudflare.UpdateCloudflareForDynamicIp(config)

	select {}
}

func loadConfig() *cloudflare.CloudflareConfig {
	return &cloudflare.CloudflareConfig{
		ZoneID: os.Getenv("CLOUDFLARE_ZONE_ID"),
		APIKey: os.Getenv("CLOUDFLARE_API_KEY"),
	}
}