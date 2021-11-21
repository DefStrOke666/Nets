package utils

import (
	"github.com/golang/freetype/truetype"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const (
	arcadeFontBaseSize = 8
	scaleMax           = 8
)

var (
	arcadeFonts map[int]font.Face
	menuFonts   map[int]font.Face
)

func getFontByPath(path string) *truetype.Font {
	file, err := os.Open("assets/MachineGunk.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	tt, err := truetype.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	return tt
}

func GetArcadeFonts(scale int) font.Face {
	if arcadeFonts == nil {
		tt := getFontByPath("assets/MachineGunk.ttf")

		arcadeFonts = map[int]font.Face{}
		for i := 1; i <= scaleMax; i++ {
			const dpi = 72
			arcadeFonts[i] = truetype.NewFace(tt, &truetype.Options{
				Size:    float64(arcadeFontBaseSize * i),
				DPI:     dpi,
				Hinting: font.HintingFull,
			})
		}
	}
	return arcadeFonts[scale]
}

func GetMenuFonts(scale int) font.Face {
	if menuFonts == nil {
		tt := getFontByPath("assets/SerpensRegular.ttf")

		menuFonts = map[int]font.Face{}
		for i := 1; i <= scaleMax; i++ {
			const dpi = 72
			menuFonts[i] = truetype.NewFace(tt, &truetype.Options{
				Size:    float64(arcadeFontBaseSize * i),
				DPI:     dpi,
				Hinting: font.HintingFull,
			})
		}
	}
	return menuFonts[scale]
}

func TextWidth(str string, fontFace font.Face) int {
	maxW := 0
	for _, line := range strings.Split(str, "\n") {
		b, _ := font.BoundString(fontFace, line)
		w := (b.Max.X - b.Min.X).Ceil()
		if maxW < w {
			maxW = w
		}
	}
	return maxW
}

func TextHeight(str string, fontFace font.Face) int {
	maxH := 0
	for _, line := range strings.Split(str, "\n") {
		b, _ := font.BoundString(fontFace, line)
		h := (b.Max.Y - b.Min.Y).Ceil()
		if maxH < h {
			maxH = h
		}
	}
	return maxH
}

var (
	shadowColor = color.NRGBA{0, 0, 0, 0x80}
)

func drawTextWithShadow(rt *ebiten.Image, str string, x, y, scale int, clr color.Color) {
	offsetY := arcadeFontBaseSize * scale
	for _, line := range strings.Split(str, "\n") {
		y += offsetY
		text.Draw(rt, line, GetArcadeFonts(scale), x+1, y+1, shadowColor)
		text.Draw(rt, line, GetArcadeFonts(scale), x, y, clr)
	}
}

func drawTextWithShadowCenter(rt *ebiten.Image, str string, x, y, scale int, clr color.Color, width int) {
	w := TextWidth(str, GetArcadeFonts(scale))
	x += (width - w) / 2
	drawTextWithShadow(rt, str, x, y, scale, clr)
}

func drawTextWithShadowRight(rt *ebiten.Image, str string, x, y, scale int, clr color.Color, width int) {
	w := TextWidth(str, GetArcadeFonts(scale))
	x += width - w
	drawTextWithShadow(rt, str, x, y, scale, clr)
}

func CreateStringImage(str string, fontFace font.Face, clr color.Color) *ebiten.Image {
	w := TextWidth(str, fontFace)
	h := TextHeight(str, fontFace)

	img := ebiten.NewImage(w, h)
	img.Clear()
	text.Draw(img, str, fontFace, 0, h, clr)

	return img
}
