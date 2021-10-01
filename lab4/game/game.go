package game

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

var (
	closeWindow  = false
	sizeChanged  = false
	screenWidth  = 720
	screenHeight = 480
)

type Game struct {
	sceneManager *SceneManager
	input        Input
}

func (g *Game) Update() error {
	if g.sceneManager == nil {
		g.sceneManager = &SceneManager{}
		g.sceneManager.GoTo(NewTitleScene())
	}

	g.input.Update()
	if err := g.sceneManager.Update(&g.input); err != nil {
		return err
	}
	if closeWindow {
		return errors.New("Closed")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
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

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
