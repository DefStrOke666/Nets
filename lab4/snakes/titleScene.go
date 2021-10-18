package snakes

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var imageBackground *ebiten.Image

func init() {
	imageBackground = getRoundRect(15, 15, backgroundColor)
}

type TitleScene struct {
	pics []*Picture

	count int
}

func NewTitleScene() *TitleScene {
	scene := &TitleScene{}

	scene.pics = make([]*Picture, 4)
	scene.pics[0] = NewPicture(
		createStringImage("SNAKES", getMenuFonts(8), titleIdleColor),
		createStringImage("SNAKES", getMenuFonts(8), titleActiveColor))
	scene.pics[1] = NewPicture(
		createStringImage("Create", getArcadeFonts(8), idleColor),
		createStringImage("Create", getArcadeFonts(8), activeColor),
	).SetHandler(func(state *GameState) {
		state.SceneManager.GoTo(NewCreateScene())
	})
	scene.pics[2] = NewPicture(
		createStringImage("Join", getArcadeFonts(8), idleColor),
		createStringImage("Join", getArcadeFonts(8), activeColor),
	).SetHandler(func(state *GameState) {
		println("server list")
		state.SceneManager.GoTo(NewJoinScene())
	})
	scene.pics[3] = NewPicture(
		createStringImage("Exit", getArcadeFonts(8), idleColor),
		createStringImage("Exit", getArcadeFonts(8), activeColor),
	).SetHandler(func(state *GameState) {
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
				s.pics[i].Handle(state)
			}
		} else {
			s.pics[i].SetActive(false)
		}
	}

	return nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(fillColor)
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
