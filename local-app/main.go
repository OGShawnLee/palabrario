package main

import (
	"context"
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's embed package to embed the frontend files into the binary.
//
//go:embed all:frontend/dist
var assets embed.FS

func main() {
	appService := NewApp()

	app := application.New(application.Options{
		Name:        "Palabrario",
		Description: "Plataforma dialectal panhispánica y arena multijugador",
		Services: []application.Service{
			application.NewService(appService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	appService.SetContext(context.Background())

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Palabrario",
		Width:            1280,
		Height:           820,
		MinWidth:         390,
		MinHeight:        680,
		BackgroundColour: application.NewRGB(246, 243, 236), // #F6F3EC (paper)
		URL:              "/",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
