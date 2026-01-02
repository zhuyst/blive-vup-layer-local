package main

import (
	"blive-vup-layer/config"
	"blive-vup-layer/util"
	"embed"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"golang.design/x/hotkey"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var iconFS embed.FS

const Name = "巫女酱子弹幕姬"

func logError(errMsg string) {
	util.ShowErrorDialog(errMsg)
	log.Error(errMsg)
}

func main() {
	defer util.Recover()

	log.SetFormatter(&log.JSONFormatter{})

	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		logError(fmt.Sprintf("failed to create logs dir: %v", err))
		return
	}
	logFile, err := os.OpenFile(fmt.Sprintf("logs/%s.txt", time.Now().Format("2006-01-02-15-04-05")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logError(fmt.Sprintf("failed to create log file: %v", err))
		return
	}
	defer logFile.Close()
	if _, err := logFile.Write([]byte("log file created\n")); err != nil {
		logError(fmt.Sprintf("failed to write log file: %v", err))
		return
	}

	logWriter := util.NewAppLogWriter(logFile)
	log.SetOutput(logWriter)

	os.RemoveAll(config.ResultFilePath)
	if err := os.MkdirAll(config.ResultFilePath, 0755); err != nil {
		log.Fatalf("os.MkdirAll err: %v", err.Error())
		return
	}
	slogLogger := slog.New(
		slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}),
	)

	service := NewService(logWriter)

	var mainWindow *application.WebviewWindow
	var app *application.App
	app = application.New(application.Options{
		Name:        Name,
		Description: Name,
		Services: []application.Service{
			application.NewServiceWithOptions(service, application.ServiceOptions{
				Route: "/result/",
			}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Windows: application.WindowsOptions{
			WebviewUserDataPath: "./tmp/",
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "com.mikocat.blive-vup-layer",
			OnSecondInstanceLaunch: func(_ application.SecondInstanceData) {
				if mainWindow != nil {
					mainWindow.Show()
					mainWindow.Restore()
					mainWindow.Focus()
				}
			},
		},
		PanicHandler: func(detail *application.PanicDetails) {
			log.Errorf("panic: %s, stack: %s", detail.Error, detail.FullStackTrace)

			dialog := app.Dialog.Error()
			dialog.SetTitle("程序发生崩溃，已恢复")
			dialog.SetMessage(fmt.Sprintf("程序发生崩溃，已恢复\npanic: %s\nstack: %s", detail.Error, detail.FullStackTrace))
			dialog.Show()
		},
		//KeyBindings: map[string]func(window *application.WebviewWindow){
		//	"F6": func(window *application.WebviewWindow) {
		//		service.writeResultOK(ResultTypeRecordStateChange, nil)
		//	},
		//	"F7": func(window *application.WebviewWindow) {
		//		service.writeResultOK(ResultTypeRecordStart, nil)
		//	},
		//	"F8": func(window *application.WebviewWindow) {
		//		service.writeResultOK(ResultTypeRecordStop, nil)
		//	},
		//},
		Logger: slogLogger,
	})
	util.SetApp(app)

	systemTray := app.SystemTray.New()
	systemTray.SetLabel(Name)
	systemTrayMenu := app.NewMenu()

	a := &App{
		App:            app,
		SystemTrayMenu: systemTrayMenu,
		WindowMap:      make(map[string]*SubWindow),
	}
	service.Init(a)

	mainWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            fmt.Sprintf("%s - 总控台", Name),
		Width:            1600,
		Height:           900,
		AlwaysOnTop:      true,
		BackgroundColour: application.NewRGB(255, 255, 255),
		URL:              "/",
	})
	mainWindow.OnWindowEvent(events.Windows.WindowClosing, func(_ *application.WindowEvent) {
		app.Quit()
	})

	addWindowMenuItem(&addWindowMenuItemParams{
		SystemTrayMenu: systemTrayMenu,
		Window:         mainWindow,
		Label:          "显示总控台",
	})
	systemTrayMenu.Add("显示所有窗口").OnClick(func(_ *application.Context) {
		if mainWindow != nil {
			mainWindow.Show()
			mainWindow.Restore()
			mainWindow.Focus()
		}
		for _, subWindow := range a.WindowMap {
			subWindow.Show()
		}
	})
	systemTrayMenu.AddSeparator()
	a.AddSubWindow(&AddSubWindowParams{
		ID:   "danmu",
		Name: "弹幕",
		URL:  "/#/danmu",
	})
	a.AddSubWindow(&AddSubWindowParams{
		ID:   "enter-room",
		Name: "进入直播间",
		URL:  "/#/enter_room",
	})
	a.AddSubWindow(&AddSubWindowParams{
		ID:   "gift",
		Name: "礼物",
		URL:  "/#/gift",
	})
	a.AddSubWindow(&AddSubWindowParams{
		ID:   "interact-word",
		Name: "关注直播间",
		URL:  "/#/interact_word",
	})
	a.AddSubWindow(&AddSubWindowParams{
		ID:   "llm",
		Name: "大模型回复",
		URL:  "/#/llm",
	})
	a.AddSubWindow(&AddSubWindowParams{
		ID:   "membership",
		Name: "大航海",
		URL:  "/#/membership",
	})
	a.AddSubWindow(&AddSubWindowParams{
		ID:   "superchat",
		Name: "醒目留言",
		URL:  "/#/superchat",
	})
	systemTrayMenu.AddSeparator()
	systemTrayMenu.Add("退出").OnClick(func(_ *application.Context) {
		app.Quit()
	})

	systemTray.SetMenu(systemTrayMenu)
	//systemTray.Run()

	recordStateChangeHk, err := registerHotKey([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModAlt}, hotkey.KeyF6, func() {
		service.writeResultOK(ResultTypeRecordStateChange, nil)
	})
	if err != nil {
		log.Fatalf("registerHotKey Ctrl+Alt+F6 err: %v", err)
		return
	}
	defer recordStateChangeHk.Unregister()

	recordStartHk, err := registerHotKey([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModAlt}, hotkey.KeyF7, func() {
		service.writeResultOK(ResultTypeRecordStart, nil)
	})
	if err != nil {
		log.Fatalf("registerHotKey Ctrl+Alt+F7 err: %v", err)
		return
	}
	defer recordStartHk.Unregister()

	recordStopHk, err := registerHotKey([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModAlt}, hotkey.KeyF8, func() {
		service.writeResultOK(ResultTypeRecordStop, nil)
	})
	if err != nil {
		log.Fatalf("registerHotKey Ctrl+Alt+F8 err: %v", err)
		return
	}
	defer recordStopHk.Unregister()

	log.Infof("App started")
	if err := app.Run(); err != nil {
		log.Fatalf("app.Run err: %v", err)
		return
	}
	log.Infof("App exited")
}

type App struct {
	App            *application.App
	SystemTrayMenu *application.Menu
	WindowMap      map[string]*SubWindow
}

type SubWindow struct {
	ID       string
	Window   *application.WebviewWindow
	MenuItem *application.MenuItem
}

func (window *SubWindow) Show() {
	window.Window.Show()
	window.Window.Restore()
	window.Window.Focus()
}

type AddSubWindowParams struct {
	ID   string
	Name string
	URL  string
}

func (app *App) AddSubWindow(params *AddSubWindowParams) *SubWindow {
	window := app.App.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: fmt.Sprintf("%s - %s", Name, params.Name),
		//Width:  1600,
		//Height: 900,
		Width:            800,
		Height:           510,
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		URL:              params.URL,
		Frameless:        true,
		DisableResize:    true,
		BackgroundType:   application.BackgroundTypeTransparent,
		Windows: application.WindowsWindow{
			DisableFramelessWindowDecorations: true,
		},
	})

	menuItemName := fmt.Sprintf("显示%s", params.Name)
	menuItem := addWindowMenuItem(&addWindowMenuItemParams{
		SystemTrayMenu: app.SystemTrayMenu,
		Window:         window,
		Label:          menuItemName,
	})

	subWindow := &SubWindow{
		Window:   window,
		MenuItem: menuItem,
	}
	app.WindowMap[params.ID] = subWindow
	return subWindow
}

type addWindowMenuItemParams struct {
	SystemTrayMenu *application.Menu
	Window         *application.WebviewWindow
	Label          string
}

func addWindowMenuItem(params *addWindowMenuItemParams) *application.MenuItem {
	window := params.Window
	windowMenuItem := params.SystemTrayMenu.AddCheckbox(params.Label, true).OnClick(func(c *application.Context) {
		if window == nil {
			return
		}
		if c.IsChecked() {
			window.Show()
			window.Restore()
			window.Focus()
		} else {
			window.Hide()
		}
	})
	window.OnWindowEvent(events.Windows.WindowHide, func(_ *application.WindowEvent) {
		windowMenuItem.SetChecked(false)
	})
	window.OnWindowEvent(events.Windows.WindowShow, func(_ *application.WindowEvent) {
		windowMenuItem.SetChecked(true)
	})
	return windowMenuItem
}

func registerHotKey(mods []hotkey.Modifier, key hotkey.Key, handler func()) (*hotkey.Hotkey, error) {
	hk := hotkey.New(mods, key)
	if err := hk.Register(); err != nil {
		return nil, err
	}

	go func() {
		for range hk.Keydown() {
			handler()
		}
	}()
	return hk, nil
}
