package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type CreateScene struct {
	backgroundPics []*utils.Picture
	buttonPics     []*utils.Picture
	background     *ebiten.Image
}

func NewCreateScene() *CreateScene {
	scene := &CreateScene{}

	scene.backgroundPics = make([]*utils.Picture, 2)
	scene.buttonPics = make([]*utils.Picture, 2)

	scene.updateImages()

	return scene
}
func (c *CreateScene) updateImages() {
	margin := int(utils.Margin)
	spacingsV := margin * 3
	spacingsH := margin * 2

	titleH := utils.TextHeight("Create game", utils.GetMenuFonts(8)) + margin

	widthUnit := (screenWidth - spacingsH) / 10
	heightUnit := (screenHeight - titleH - spacingsV) / 6

	configW := widthUnit * 10
	configH := heightUnit * 5

	c.background = utils.GetRoundRect(screenWidth, screenHeight, utils.BackgroundColor)
	c.backgroundPics[0] = utils.NewPicture(
		utils.CreateStringImage("Create game", utils.GetMenuFonts(8), utils.TitleIdleColor),
		utils.CreateStringImage("Create game", utils.GetMenuFonts(8), utils.TitleActiveColor))
	c.backgroundPics[1] = utils.NewPicture(
		utils.GetRoundRectWithBorder(configW, configH, utils.CentreIdleColor, utils.LineIdleColor),
		utils.GetRoundRectWithBorder(configW, configH, utils.CentreActiveColor, utils.LineActiveColor))

	c.backgroundPics[0].SetRect(c.backgroundPics[0].GetIdleImage().Bounds().Add(image.Point{X: margin, Y: margin}))
	c.backgroundPics[1].SetRect(c.backgroundPics[1].GetIdleImage().Bounds().Add(image.Point{X: margin, Y: titleH + margin}))

	buttonW := widthUnit * 3
	buttonH := heightUnit

	c.buttonPics[0] = utils.NewPicture(
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreIdleColor, utils.LineIdleColor, "Start", utils.GetMenuFonts(4)),
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreActiveColor, utils.LineActiveColor, "Start", utils.GetMenuFonts(4)),
	).SetHandler(func() {
		sceneManager.GoTo(NewGameScene(utils.NewDefaultGameConfig()))
	})
	c.buttonPics[1] = utils.NewPicture(
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreIdleColor, utils.LineIdleColor, "Return", utils.GetMenuFonts(4)),
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreActiveColor, utils.LineActiveColor, "Return", utils.GetMenuFonts(4)),
	).SetHandler(func() {
		sceneManager.GoTo(NewTitleScene())
	})

	c.buttonPics[0].SetRect(c.buttonPics[0].GetIdleImage().Bounds().Add(image.Point{X: margin, Y: titleH + margin*2 + configH}))
	c.buttonPics[1].SetRect(c.buttonPics[1].GetIdleImage().Bounds().Add(image.Point{X: screenWidth - margin - buttonW, Y: titleH + margin*2 + configH}))
}

func (c *CreateScene) Update(state *GameState) error {
	if sizeChanged {
		c.updateImages()
	}

	for i := range c.buttonPics {
		c.buttonPics[i].Update()
	}

	for i := range c.backgroundPics {
		c.backgroundPics[i].Update()
	}

	return nil
}

func (c *CreateScene) Draw(screen *ebiten.Image) {
	screen.Fill(utils.FillColor)
	screen.DrawImage(c.background, &ebiten.DrawImageOptions{})

	for i := range c.backgroundPics {
		c.backgroundPics[i].Draw(screen)
	}

	for i := range c.buttonPics {
		c.buttonPics[i].Draw(screen)
	}
}
