package main

import (
	"blive-vup-layer/config"
	"embed"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

//go:embed all:frontend/dist
var assets embed.FS

type ResultFileLoader struct {
	http.Handler
}

func (h *ResultFileLoader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/result/") {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fileName := strings.TrimPrefix(r.URL.Path, "/result/")
	filePath := path.Join(config.ResultFilePath, fileName)
	f, err := os.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Errorf("fileName: %s, filePath: %s not found", fileName, filePath)
		return
	}

	w.Write(f)
}

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

	os.RemoveAll(config.ResultFilePath)
	if err := os.MkdirAll(config.ResultFilePath, 0755); err != nil {
		log.Fatalf("os.MkdirAll err: %v", err.Error())
		return
	}

	app := NewApp(logWriter)
	if err := wails.Run(&options.App{
		Title:  "巫女酱子弹幕姬",
		Width:  1600,
		Height: 900,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: &ResultFileLoader{},
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		//Frameless: true,
	}); err != nil {
		log.Errorf("wails.Run err: %v", err)
	}
}
