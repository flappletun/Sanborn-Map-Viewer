package main

import (
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// base file paths
var basePath = "/Users/samwalker/Library/Mobile Documents/com~apple~CloudDocs/Projects/Go/Sanborn Map Explorer/resources"
var baseThumbnailPath = filepath.Join(basePath, "thumbnails")
var baseMediumPath = filepath.Join(basePath, "small_size")
var baseLargePath = filepath.Join(basePath, "full_size")
var gopherPath = filepath.Join(basePath, "gopher.png")

var townList = []string{
	"Hallettsville",
	"Moulton",
	"Shiner",
	"Yoakum",
}

var imageMap = NewLockedImageMap()
var wg sync.WaitGroup

func main() {
	var thumbnails []fyne.CanvasObject
	var names []string
	var thumbnailGrid *container.Scroll

	// initial thumbnail grid with gopher image as placeholder
	gopher := canvas.NewImageFromFile(gopherPath)
	gopher.FillMode = canvas.ImageFillContain
	gopher.SetMinSize(fyne.NewSize(800, 800))
	thumbnailGrid = container.NewScroll(gopher)
	thumbnailGrid.SetMinSize(fyne.NewSize(630, 630))

	// drop down menu to select town
	townSelector := widget.NewSelect(townList, func(s string) {
		// load thumbnails for selected town
		thumbnails, names = LoadThumbnails(s)
		grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 200)), thumbnails...)
		thumbnailGrid.Content = grid
		thumbnailGrid.Refresh()

		// load small images for each thumbnail
		wg.Add(len(names))
		for _, name := range names {
			go LoadMapWorker(name)
		}
		wg.Wait()
	})
	townSelector.PlaceHolder = "Select a town"

	// run app
	myApp := app.New()
	myWindow := myApp.NewWindow("Sanborn Map Explorer")
	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.SetContent(container.NewHBox(townSelector, thumbnailGrid))
	myWindow.Show()
	myApp.Run()
	tidyUp()
}

func tidyUp() {}
