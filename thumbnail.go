package main

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

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

// loads thumbnails of each avaiable map for the selected town
func LoadThumbnails(townName string) ([]fyne.CanvasObject, []string) {
	thumbnails := make([]fyne.CanvasObject, 0, 64)
	names := make([]string, 0, 64)
	thumbnailDir := filepath.Join(baseThumbnailPath, townName)
	err := filepath.Walk(thumbnailDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing path:", err)
		}
		if filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".png" {
			thumbnails = append(thumbnails, NewThumbnail(path, func() {
				println("Thumbnail clicked:", filepath.Base(path))

				// load medium image corresponding to the thumbnail
				//todo: use image map
			}))
			names = append(names, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error loading thumbnails:", err)
	}
	return thumbnails, names
}
