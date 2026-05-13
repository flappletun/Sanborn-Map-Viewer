package main

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var defaultCapacity = 30

var townList = []string{
	"Hallettsville",
	"Moulton",
	"Shiner",
	"Yoakum",
}

type Thumbnail struct {
	widget.BaseWidget
	Image   *canvas.Image
	OnClick func()
	Name    string
}

func NewThumbnail(path string, onClick func()) *Thumbnail {
	img := canvas.NewImageFromFile(path)
	img.FillMode = canvas.ImageFillContain
	t := &Thumbnail{
		Image:   img,
		OnClick: onClick,
		Name:    filepath.Base(path),
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *Thumbnail) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.Image)
}

func (t *Thumbnail) Tapped(_ *fyne.PointEvent) {
	if t.OnClick != nil {
		t.OnClick()
	}
}

// base file paths
var basePath = "/Users/samwalker/Library/Mobile Documents/com~apple~CloudDocs/Projects/Go/Sanborn Map Explorer/resources"
var baseThumbnailPath = filepath.Join(basePath, "thumbnails")
var gopherPath = filepath.Join(basePath, "gopher.png")

func loadThumbnails(townName string) []fyne.CanvasObject {
	thumbnails := make([]fyne.CanvasObject, 0, defaultCapacity)
	thumbnailDir := filepath.Join(baseThumbnailPath, townName)
	err := filepath.Walk(thumbnailDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing path:", err)
		}
		if filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".png" {
			thumbnails = append(thumbnails, NewThumbnail(path, func() {
				println("Thumbnail clicked:", filepath.Base(path))
			}))
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error loading thumbnails:", err)
	}
	return thumbnails
}

func main() {
	var thumbnails []fyne.CanvasObject
	var thumbnailGrid *container.Scroll

	// initial thumbnail grid with gopher image as placeholder
	gopher := canvas.NewImageFromFile(gopherPath)
	gopher.FillMode = canvas.ImageFillContain
	gopher.SetMinSize(fyne.NewSize(800, 800))
	thumbnailGrid = container.NewScroll(gopher)
	thumbnailGrid.SetMinSize(fyne.NewSize(630, 630))

	// drop down menu to select town
	townSelector := widget.NewSelect(townList, func(s string) {
		thumbnails = loadThumbnails(s)
		grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 200)), thumbnails...)
		thumbnailGrid.Content = grid
		thumbnailGrid.Refresh()
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

func tidyUp() {
	fmt.Println("Exited")
}
