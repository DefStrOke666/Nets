package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var imageBackground *ebiten.Image

func init() {
	imageBackground = utils.GetRoundRect(15, 15, utils.BackgroundColor)
}

type TitleScene struct {
	pics []*utils.Picture

	count int
}

func NewTitleScene() *TitleScene {
	scene := &TitleScene{}

	scene.pics = make([]*utils.Picture, 4)
	scene.pics[0] = utils.NewPicture(
		utils.CreateStringImage("SNAKES", utils.GetMenuFonts(8), utils.TitleIdleColor),
		utils.CreateStringImage("SNAKES", utils.GetMenuFonts(8), utils.TitleActiveColor))
	scene.pics[1] = utils.NewPicture(
		utils.CreateStringImage("Create", utils.GetArcadeFonts(8), utils.MenuTextIdleColor),
		utils.CreateStringImage("Create", utils.GetArcadeFonts(8), utils.MenuTextActiveColor),
	).SetHandler(func() {
		sceneManager.GoTo(NewCreateScene())
	})
	scene.pics[2] = utils.NewPicture(
		utils.CreateStringImage("Join", utils.GetArcadeFonts(8), utils.MenuTextIdleColor),
		utils.CreateStringImage("Join", utils.GetArcadeFonts(8), utils.MenuTextActiveColor),
	).SetHandler(func() {
		println("server list")
		sceneManager.GoTo(NewJoinScene())
	})
	scene.pics[3] = utils.NewPicture(
		utils.CreateStringImage("Exit", utils.GetArcadeFonts(8), utils.MenuTextIdleColor),
		utils.CreateStringImage("Exit", utils.GetArcadeFonts(8), utils.MenuTextActiveColor),
	).SetHandler(func() {
		println("exit")
		closeWindow = true
	})

	scene.updateImgs()

	return scene
}

func (s *TitleScene) updateImgs() {
	margin := 50
	for i := range s.pics {
		w, h := s.pics[i].GetIdleImage().Size()
		if i == 0 {
			s.pics[i].SetRect(s.pics[i].GetIdleImage().Bounds().Add(image.Point{X: (screenWidth - w) / 2, Y: h}))
		} else {
			s.pics[i].SetRect(s.pics[i].GetIdleImage().Bounds().Add(image.Point{X: (screenWidth - w) / 2, Y: margin*(i+1) + h}))
		}
	}
}

func (s *TitleScene) Update(state *GameState) error {
	if sizeChanged {
		s.updateImgs()
	}

	s.count++
	for i := range s.pics {
		if s.pics[i].InBounds(ebiten.CursorPosition()) {
			s.pics[i].SetActive(true)
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				s.pics[i].Handle()
			}
		} else {
			s.pics[i].SetActive(false)
		}
	}

	return nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(utils.FillColor)
	s.drawTitleBackground(screen, s.count)
	for i := range s.pics {
		s.pics[i].Draw(screen)
	}
}

func (s *TitleScene) drawTitleBackground(screen *ebiten.Image, c int) {
	w, h := imageBackground.Size()
	op := &ebiten.DrawImageOptions{}
	for i := 0; i < (screenWidth/w+1)*(screenHeight/h+2); i++ {
		op.GeoM.Reset()
		dx := -(c / 4) % w
		dy := (c / 4) % h
		dstX := (i%(screenWidth/w+1))*w + dx
		dstY := (i/(screenWidth/w+1)-1)*h + dy
		op.GeoM.Translate(float64(dstX), float64(dstY))
		screen.DrawImage(imageBackground, op)
	}
}
