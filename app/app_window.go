package main

import "github.com/wailsapp/wails/v3/pkg/application"

func (a *App) OnMinimize() {
	if app := application.Get(); app != nil {
		if w := app.Window.Current(); w != nil {
			w.Minimise()
		}
	}
}

func (a *App) OnMaximize() {
	if app := application.Get(); app != nil {
		if w := app.Window.Current(); w != nil {
			w.ToggleMaximise()
		}
	}
}

func (a *App) OnClose() {
	if app := application.Get(); app != nil {
		app.Quit()
	}
}
