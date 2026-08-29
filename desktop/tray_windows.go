//go:build windows

package main

import (
	_ "embed"
	"sync"

	"github.com/getlantern/systray"
)

var systemTrayOnce sync.Once

//go:embed build/windows/icon.ico
var systemTrayIcon []byte

func startSystemTray(app *App) {
	systemTrayOnce.Do(func() {
		go systray.Run(func() {
			systray.SetIcon(systemTrayIcon)
			systray.SetTitle("Lingma Proxy")
			systray.SetTooltip("Lingma Proxy - 后台运行中")

			showItem := systray.AddMenuItem("打开主窗口", "显示 Lingma Proxy 主窗口")
			systray.AddSeparator()
			quitItem := systray.AddMenuItem("完全退出", "停止代理并退出 Lingma Proxy")

			go func() {
				for {
					select {
					case <-showItem.ClickedCh:
						app.ShowWindow()
					case <-quitItem.ClickedCh:
						systray.Quit()
						app.ForceQuitApp()
						return
					}
				}
			}()
		}, func() {})
	})
}
