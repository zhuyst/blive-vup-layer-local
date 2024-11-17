package main

import (
	"blive-vup-layer/config"
	"embed"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"io"
	"os"
	"time"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatalf("failed to create logs dir: %v", err)
		return
	}
	logFile, err := os.OpenFile(fmt.Sprintf("logs/%s.txt", time.Now().Format("2006-01-02-15-04-05")), os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create log file: %v", err)
		return
	}
	logWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(logWriter)

	configFilePath := flag.String("config", "./etc/config-dev.toml", "config file path")
	flag.Parse()
	cfg, err := config.ParseConfig(*configFilePath)
	if err != nil {
		log.Fatalf("failed to parse config file: %v", err)
		return
	}

	os.RemoveAll(config.ResultFilePath)
	if err := os.MkdirAll(config.ResultFilePath, 0755); err != nil {
		log.Fatalf("os.MkdirAll err: %v", err.Error())
		return
	}

	app, err := NewApp(cfg, logWriter)
	if err != nil {
		log.Fatalf("NewApp err: %v", err.Error())
		return
	}

	if err := wails.Run(&options.App{
		Title:  "巫女酱子弹幕姬",
		Width:  1920,
		Height: 1080,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	}); err != nil {
		log.Errorf("wails.Run err: %v", err)
	}
}
