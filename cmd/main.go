package main

import (
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/maxmorhardt/cloudflare-k8s-manager/internal/k8s"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
  	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{                                             
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {                                                     
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)       
			return "", fileName                                                      
		},                                                                           
	})

	log.Info("Starting watcher")
	go k8s.Watcher()

	select {}
}

