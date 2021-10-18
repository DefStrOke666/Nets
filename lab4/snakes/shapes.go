package snakes

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"image/color"
)

func getRoundRectWithTextAndBorder(w, h int, clr color.Color, lineClr color.Color, str string, fontfFace font.Face, txtClr color.Color) *ebiten.Image {
	dc := gg.NewContext(w, h)
	dc.DrawRoundedRectangle(0, 0, float64(w), float64(h), radius)
	dc.SetRGBA(colorToScale(lineClr))
	dc.Fill()
	dc.DrawRoundedRectangle(0+lineThickness, 0+lineThickness, float64(w)-lineThickness*2, float64(h)-lineThickness*2, radius)
	dc.SetRGBA(colorToScale(clr))
	dc.Fill()

	strImg := createStringImage(str, fontfFace, txtClr)
	x, y := strImg.Size()
	centreW := (w - x) / 2
	centreH := (h - y) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(centreW), float64(centreH))
	img := ebiten.NewImageFromImage(dc.Image())
	img.DrawImage(strImg, op)

	return img
}

func getRoundRectWithBorder(w, h int, clr color.Color, lineClr color.Color) *ebiten.Image {
	dc := gg.NewContext(w, h)
	dc.DrawRoundedRectangle(0, 0, float64(w), float64(h), radius)
	dc.SetRGBA(colorToScale(lineClr))
	dc.Fill()
	dc.DrawRoundedRectangle(0+lineThickness, 0+lineThickness, float64(w)-lineThickness*2, float64(h)-lineThickness*2, radius)
	dc.SetRGBA(colorToScale(clr))
	dc.Fill()
	return ebiten.NewImageFromImage(dc.Image())
}

func getRoundRect(w, h int, clr color.Color) *ebiten.Image {
	dc := gg.NewContext(w, h)
	dc.DrawRoundedRectangle(0, 0, float64(w), float64(h), radius)
	dc.SetRGBA(colorToScale(clr))
	dc.Fill()
	return ebiten.NewImageFromImage(dc.Image())
}
