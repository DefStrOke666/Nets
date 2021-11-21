package utils

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"image/color"
)

func GetRectWithBorder(w, h int, clr color.Color, lineClr color.Color) *ebiten.Image {
	dc := gg.NewContext(w, h)
	dc.DrawRectangle(0, 0, float64(w), float64(h))
	dc.SetRGBA(colorToScale(lineClr))
	dc.Fill()
	dc.DrawRectangle(0+LineThickness, 0+LineThickness, float64(w)-LineThickness*2, float64(h)-LineThickness*2)
	dc.SetRGBA(colorToScale(clr))
	dc.Fill()
	return ebiten.NewImageFromImage(dc.Image())
}

func GetRoundRectWithBorder(w, h int, clr color.Color, lineClr color.Color) *ebiten.Image {
	dc := gg.NewContext(w, h)
	dc.DrawRoundedRectangle(0, 0, float64(w), float64(h), Radius)
	dc.SetRGBA(colorToScale(lineClr))
	dc.Fill()
	dc.DrawRoundedRectangle(0+LineThickness, 0+LineThickness, float64(w)-LineThickness*2, float64(h)-LineThickness*2, Radius)
	dc.SetRGBA(colorToScale(clr))
	dc.Fill()
	return ebiten.NewImageFromImage(dc.Image())
}

func BorderedRoundRectWithText(w, h int, clr color.Color, lineClr color.Color, str string, font font.Face) *ebiten.Image {
	textImg := CreateStringImage(str, font, lineClr)
	rectImg := GetRoundRectWithBorder(w, h, clr, lineClr)
	textW, textH := textImg.Size()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((w-textW)/2), float64((h-textH)/2))
	rectImg.DrawImage(textImg, op)
	return rectImg
}

func GetRoundRect(w, h int, clr color.Color) *ebiten.Image {
	dc := gg.NewContext(w, h)
	dc.DrawRoundedRectangle(0, 0, float64(w), float64(h), Radius)
	dc.SetRGBA(colorToScale(clr))
	dc.Fill()
	return ebiten.NewImageFromImage(dc.Image())
}
