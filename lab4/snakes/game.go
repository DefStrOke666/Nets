package snakes

import (
	"errors"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

var (
	closeWindow  = false
	sizeChanged  = false
	screenWidth  = 1100
	screenHeight = 720
	sceneManager *SceneManager
)

type Game struct{}

func (g *Game) Update() error {
	if sceneManager == nil {
		sceneManager = &SceneManager{}
		sceneManager.GoTo(NewTitleScene())
	}

	if err := sceneManager.Update(); err != nil {
		return err
	}
	if closeWindow {
		return errors.New("closed")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	sceneManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != screenWidth || outsideHeight != screenHeight {
		screenWidth = outsideWidth
		screenHeight = outsideHeight
		sizeChanged = true
	} else {
		sizeChanged = false
	}
	return screenWidth, screenHeight
}

func Play() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snakes")
	ebiten.SetWindowResizable(true)

	img, err := gg.LoadImage("assets/snake_icon.png")
	if err != nil {
		log.Fatal(err)
	}
	imgs := make([]image.Image, 1)
	imgs[0] = img
	ebiten.SetWindowIcon(imgs)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
