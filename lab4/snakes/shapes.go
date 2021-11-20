package snakes

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"image/color"
)

func getRectWithBorder(w, h int, clr color.Color, lineClr color.Color) *ebiten.Image {
	dc := gg.NewContext(w, h)
	dc.DrawRectangle(0, 0, float64(w), float64(h))
	dc.SetRGBA(colorToScale(lineClr))
	dc.Fill()
	dc.DrawRectangle(0+lineThickness, 0+lineThickness, float64(w)-lineThickness*2, float64(h)-lineThickness*2)
	dc.SetRGBA(colorToScale(clr))
	dc.Fill()
	return ebiten.NewImageFromImage(dc.Image())
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

func borderedRoundRectWithText(w, h int, clr color.Color, lineClr color.Color, str string, font font.Face) *ebiten.Image {
	textImg := createStringImage(str, font, lineClr)
	rectImg := getRoundRectWithBorder(w, h, clr, lineClr)
	textW, textH := textImg.Size()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((w-textW)/2), float64((h-textH)/2))
	rectImg.DrawImage(textImg, op)
	return rectImg
}

func getRoundRect(w, h int, clr color.Color) *ebiten.Image {
	dc := gg.NewContext(w, h)
	dc.DrawRoundedRectangle(0, 0, float64(w), float64(h), radius)
	dc.SetRGBA(colorToScale(clr))
	dc.Fill()
	return ebiten.NewImageFromImage(dc.Image())
}
