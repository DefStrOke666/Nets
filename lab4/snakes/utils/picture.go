package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

type Picture struct {
	idleImage   *ebiten.Image
	activeImage *ebiten.Image
	rect        image.Rectangle
	active      bool
	handler     func()
}

func NewPicture(idleImage *ebiten.Image, activeImage *ebiten.Image) *Picture {
	pic := &Picture{}

	pic.idleImage = idleImage
	pic.activeImage = activeImage
	pic.active = false
	pic.handler = func() {
		return
	}

	return pic
}

func (p *Picture) SetIdleImage(image *ebiten.Image) {
	p.idleImage = image
}

func (p *Picture) GetIdleImage() *ebiten.Image {
	return p.idleImage
}

func (p *Picture) SetActiveImage(image *ebiten.Image) {
	p.activeImage = image
}

func (p *Picture) GetActiveImage() *ebiten.Image {
	return p.activeImage
}

func (p *Picture) SetRect(rect image.Rectangle) {
	p.rect = rect
}

func (p *Picture) GetRect() image.Rectangle {
	return p.rect
}

func (p *Picture) SetActive(val bool) {
	p.active = val
}

func (p *Picture) IsActive() bool {
	return p.active
}

func (p *Picture) SetHandler(handler func()) *Picture {
	p.handler = handler
	return p
}

func (p *Picture) Handle() {
	p.handler()
}

func (p *Picture) InBounds(x, y int) bool {
	if x > p.rect.Min.X && x < p.rect.Max.X && y > p.rect.Min.Y && y < p.rect.Max.Y {
		return true
	}
	return false
}

func (p *Picture) Update() {
	if p.InBounds(ebiten.CursorPosition()) {
		p.SetActive(true)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			p.Handle()
		}
	} else {
		p.SetActive(false)
	}
}

func (p *Picture) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.rect.Min.X), float64(p.rect.Min.Y))
	if p.active {
		screen.DrawImage(p.activeImage, op)
	} else {
		screen.DrawImage(p.idleImage, op)
	}
}
