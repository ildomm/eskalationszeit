package utils

import (
	"github.com/ildomm/eskalationszeit/preisgenerator/config"
	"github.com/natefinch/lumberjack"
	"go/build"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func SetupLogger(filename string) {
	var basePath = config.App.Logger.BasePath

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	basePath = gopath + config.App.Logger.BasePath

	filePath := basePath + filename
	log.Println("Logger target " + filePath)
	rotate := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    config.App.Logger.MaxSize, // megabytes
		MaxBackups: config.App.Logger.MaxBackups,
		MaxAge:     config.App.Logger.MaxAge, //days
		Compress:   config.App.Logger.Compress,
	}

	// Force system log to send to both directions
	// https://www.ardanlabs.com/blog/2013/11/using-log-package-in-go.html
	multi := io.MultiWriter(rotate, os.Stdout)
	log.SetOutput(multi)
}

func SignalNotify() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
}

func cleanup() {
	log.Println("Shutting down...")
}