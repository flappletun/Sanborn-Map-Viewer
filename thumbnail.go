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
			name := filepath.Base(path)
			thumbnails = append(thumbnails, NewThumbnail(path, func() {
				// load medium image corresponding to the thumbnail
				sanbornHolder = imageMap.Map[name]
				thumbnailFrame.Content = sanbornHolder.mediumSize
				thumbnailFrame.Refresh()

				// reveal back button to return to thumbnail grid
				backButton.Show()

				// reveal detail button to show large image
				detailButton.Show()
			}))
			names = append(names, name)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error loading thumbnails:", err)
	}
	return thumbnails, names
}
