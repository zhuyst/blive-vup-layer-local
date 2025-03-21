package main

import (
	"blive-vup-layer/config"
	"embed"
	_ "embed"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v3/pkg/application"
	"io"
	"os"
	"runtime/debug"
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

	os.RemoveAll(config.ResultFilePath)
	if err := os.MkdirAll(config.ResultFilePath, 0755); err != nil {
		log.Fatalf("os.MkdirAll err: %v", err.Error())
		return
	}

	service := NewService(logWriter)

	var mainWindow *application.WebviewWindow
	app := application.New(application.Options{
		Name:        "巫女酱子弹幕姬",
		Description: "巫女酱子弹幕姬",
		Services: []application.Service{
			application.NewService(service, application.ServiceOptions{
				Route: "/result/",
			}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Windows: application.WindowsOptions{
			WebviewUserDataPath: "./tmp/",
		},
		OnShutdown: func() {
			service.StopConn()
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "com.mikocat.blive-vup-layer",
			OnSecondInstanceLaunch: func(_ application.SecondInstanceData) {
				if mainWindow != nil {
					mainWindow.Restore()
					mainWindow.Focus()
				}
			},
		},
		PanicHandler: func(err any) {
			log.Errorf("panic: %s, stack: %s", err, string(debug.Stack()))
		},
	})
	service.startup(app)

	mainWindow = app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:            "巫女酱子弹幕姬-总控台",
		Width:            1600,
		Height:           900,
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
