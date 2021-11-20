package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type CreateScene struct {
	backgroundPics []*Picture
	buttonPics     []*Picture
	background     *ebiten.Image
}

func NewCreateScene() *CreateScene {
	scene := &CreateScene{}

	scene.backgroundPics = make([]*Picture, 2)
	scene.buttonPics = make([]*Picture, 2)

	scene.updateImages()

	return scene
}
func (c *CreateScene) updateImages() {
	margin := int(margin)
	spacingsV := margin * 3
	spacingsH := margin * 2

	titleH := textHeight("Create game", getMenuFonts(8)) + margin

	widthUnit := (screenWidth - spacingsH) / 10
	heightUnit := (screenHeight - titleH - spacingsV) / 6

	configW := widthUnit * 10
	configH := heightUnit * 5

	c.background = getRoundRect(screenWidth, screenHeight, backgroundColor)
	c.backgroundPics[0] = NewPicture(
		createStringImage("Create game", getMenuFonts(8), titleIdleColor),
		createStringImage("Create game", getMenuFonts(8), titleActiveColor))
	c.backgroundPics[1] = NewPicture(
		getRoundRectWithBorder(configW, configH, centreIdleColor, lineIdleColor),
		getRoundRectWithBorder(configW, configH, centreActiveColor, lineActiveColor))

	c.backgroundPics[0].SetRect(c.backgroundPics[0].GetIdleImage().Bounds().Add(image.Point{X: margin, Y: margin}))
	c.backgroundPics[1].SetRect(c.backgroundPics[1].GetIdleImage().Bounds().Add(image.Point{X: margin, Y: titleH + margin}))

	buttonW := widthUnit * 3
	buttonH := heightUnit

	c.buttonPics[0] = NewPicture(
		borderedRoundRectWithText(buttonW, buttonH, centreIdleColor, lineIdleColor, "Start", getMenuFonts(4)),
		borderedRoundRectWithText(buttonW, buttonH, centreActiveColor, lineActiveColor, "Start", getMenuFonts(4)),
	).SetHandler(func(s *GameState) {
		s.SceneManager.GoTo(NewGameScene(utils.NewDefaultGameConfig()))
	})
	c.buttonPics[1] = NewPicture(
		borderedRoundRectWithText(buttonW, buttonH, centreIdleColor, lineIdleColor, "Return", getMenuFonts(4)),
		borderedRoundRectWithText(buttonW, buttonH, centreActiveColor, lineActiveColor, "Return", getMenuFonts(4)),
	).SetHandler(func(s *GameState) {
		s.SceneManager.GoTo(NewTitleScene())
	})

	c.buttonPics[0].SetRect(c.buttonPics[0].GetIdleImage().Bounds().Add(image.Point{X: margin, Y: titleH + margin*2 + configH}))
	c.buttonPics[1].SetRect(c.buttonPics[1].GetIdleImage().Bounds().Add(image.Point{X: screenWidth - margin - buttonW, Y: titleH + margin*2 + configH}))
}

func (c *CreateScene) Update(state *GameState) error {
	if sizeChanged {
		c.updateImages()
	}

	for i := range c.buttonPics {
		c.buttonPics[i].Update(state)
	}

	for i := range c.backgroundPics {
		c.backgroundPics[i].Update(state)
	}

	return nil
}

func (c *CreateScene) Draw(screen *ebiten.Image) {
	screen.Fill(fillColor)
	screen.DrawImage(c.background, &ebiten.DrawImageOptions{})

	for i := range c.backgroundPics {
		c.backgroundPics[i].Draw(screen)
	}

	for i := range c.buttonPics {
		c.buttonPics[i].Draw(screen)
	}
}
