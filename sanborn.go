package main

import (
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2/canvas"
)

type LockedImageMap struct {
	Map  map[string]*canvas.Image
	Lock sync.RWMutex
}

func NewLockedImageMap() *LockedImageMap {
	return &LockedImageMap{
		Map:  make(map[string]*canvas.Image),
		Lock: sync.RWMutex{},
	}
}

type SanbornMap struct {
	Name       string
	mediumSize *canvas.Image
	largeSize  *canvas.Image
}

func NewSanbornMap(name string) *SanbornMap {
	return &SanbornMap{
		Name: name,
	}
}

func LoadMapWorker(name, town string) {
	// load medium images corrsesponing to the thumbnails
	imagePath := filepath.Join(baseMediumPath, town, name)
	println("Loading image:", imagePath)
	image := canvas.NewImageFromFile(imagePath)
	image.FillMode = canvas.ImageFillContain
	imageMap.Lock.Lock()
	imageMap.Map[name] = image
	imageMap.Lock.Unlock()
	wg.Done()
}
