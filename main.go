package main

import (
	"blive-vup-layer/config"
	"embed"
	_ "embed"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"io"
	"os"
	"runtime/debug"
	"time"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var iconFS embed.FS

const Name = "巫女酱子弹幕姬"

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
		Name:        Name,
		Description: Name,
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
		PanicHandler: func(err any) {
			log.Errorf("panic: %s, stack: %s", err, string(debug.Stack()))
		},
		KeyBindings: map[string]func(window *application.WebviewWindow){
			"F6": func(window *application.WebviewWindow) {
				service.writeResultOK(ResultTypeRecordStateChange, nil)
			},
			"F7": func(window *application.WebviewWindow) {
				service.writeResultOK(ResultTypeRecordStart, nil)
			},
			"F8": func(window *application.WebviewWindow) {
				service.writeResultOK(ResultTypeRecordStop, nil)
			},
		},
	})

	systemTray := app.NewSystemTray()
	systemTray.SetLabel(Name)
	systemTrayMenu := app.NewMenu()

	a := &App{
		App:            app,
		SystemTrayMenu: systemTrayMenu,
		WindowMap:      make(map[string]*SubWindow),
	}
	service.Init(a)

	mainWindow = app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
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

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
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
	window := app.App.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title:            fmt.Sprintf("%s - %s", Name, params.Name),
		Width:            1600,
		Height:           900,
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
