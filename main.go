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
var backButton *widget.Button
var detailButton *widget.Button
var fullViewButton *widget.Button
var sanbornHolder SanbornMap
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

		backButton.Hide()
		detailButton.Hide()
		fullViewButton.Hide()
	})
	townSelector.PlaceHolder = "Select a town"

	// back button to return to thumbnail view
	backButton = widget.NewButton("Back", func() {
		thumbnailFrame.Content = fullGrid
		thumbnailFrame.Refresh()
		backButton.Hide()
		detailButton.Hide()
		fullViewButton.Hide()
	})
	backButton.Hide() // hide back button until a thumbnail is clicked

	// detail button to show large image
	detailButton = widget.NewButton("Detailed View", func() {
		thumbnailFrame.Content = gopher //todo
		thumbnailFrame.Refresh()
		detailButton.Hide()
		backButton.Hide()
		fullViewButton.Show()
	})
	detailButton.Hide() // hide detail button until a thumbnail is clicked

	// full view button to retun to full view
	fullViewButton = widget.NewButton("Full View", func() {
		thumbnailFrame.Content = sanbornHolder.mediumSize //todo load
		thumbnailFrame.Refresh()
		detailButton.Show()
		backButton.Show()
		fullViewButton.Hide()
	})
	fullViewButton.Hide()

	// control bar
	controlBar := container.NewVBox(townSelector, backButton, detailButton, fullViewButton)

	// run app
	myApp := app.New()
	myWindow := myApp.NewWindow("Sanborn Map Explorer")
	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.SetContent(container.NewHBox(controlBar, thumbnailFrame))
	myWindow.Show()
	myApp.Run()
	tidyUp()
}

func tidyUp() {}
