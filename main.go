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

// global variables
var thumbnailFrame *container.Scroll
var imageMap = NewLockedImageMap()
var wg sync.WaitGroup

func main() {
	var thumbnails []fyne.CanvasObject
	var names []string
	var fullGrid *fyne.Container

	// initialize thumbnail grid with gopher image as placeholder
	gopher := canvas.NewImageFromFile(gopherPath)
	gopher.FillMode = canvas.ImageFillContain
	gopher.SetMinSize(fyne.NewSize(800, 800))
	thumbnailFrame = container.NewScroll(gopher)
	thumbnailFrame.SetMinSize(fyne.NewSize(630, 630))

	// drop down menu to select town
	townSelector := widget.NewSelect(townList, func(s string) {
		// load thumbnails for selected town
		thumbnails, names = LoadThumbnails(s)
		fullGrid = container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 200)), thumbnails...)
		thumbnailFrame.Content = fullGrid
		thumbnailFrame.Refresh()

		// load medium image for each thumbnail
		wg.Add(len(names))
		for _, name := range names {
			go LoadMapWorker(name, s)
		}
		wg.Wait()
	})
	townSelector.PlaceHolder = "Select a town"

	// run app
	myApp := app.New()
	myWindow := myApp.NewWindow("Sanborn Map Explorer")
	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.SetContent(container.NewHBox(townSelector, thumbnailFrame))
	myWindow.Show()
	myApp.Run()
	tidyUp()
}

func tidyUp() {}
