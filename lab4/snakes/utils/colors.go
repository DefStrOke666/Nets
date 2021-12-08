package utils

import (
	"fmt"
	"image/color"
)

var (
	FillColor = ParseHexColor("FCF8E8")

	BackgroundColor = ParseHexColor("F9F3DF")

	CentreIdleColor   = ParseHexColor("FFF4E4")
	CentreActiveColor = ParseHexColor("FFF4E4")
	LineIdleColor     = ParseHexColor("AACDBE")
	LineActiveColor   = ParseHexColor("FF7657")

	TitleIdleColor      = ParseHexColor("A2D0C1")
	TitleActiveColor    = ParseHexColor("EB5033")
	MenuTextIdleColor   = ParseHexColor("2C2E43")
	MenuTextActiveColor = ParseHexColor("36C186")

	ServerBackgroundIdleColor   = ParseHexColor("FFF4E4")
	ServerBackgroundActiveColor = ParseHexColor("FFF4E4")
	ServerTextIdleColor         = ParseHexColor("F69E7B")
	ServerTextActiveColor       = ParseHexColor("FA7D09")

	ScoreCentreColor = ParseHexColor("F6EDCF")
	ScoreLineColor   = ParseHexColor("FDD365")
	ScoreTextColor   = ParseHexColor("EB4D55")

	ConfigCentreColor = ParseHexColor("F6EDCF")
	ConfigLineColor   = ParseHexColor("FDD365")
	ConfigTextColor   = ParseHexColor("EB4D55")

	FieldCellColor1 = ParseHexColor("BFD8B8")
	FieldCellColor2 = ParseHexColor("E7EAB5")

	FoodColor       = ParseHexColor("D34848")
	SnakeBodyColor1 = ParseHexColor("FEA82F")
	SnakeHeadColor1 = ParseHexColor("FF6701")
)

func ParseHexColor(s string) color.RGBA {
	c := color.RGBA{}
	c.A = 0xff
	switch len(s) {
	case 6:
		_, _ = fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	case 3:
		_, _ = fmt.Sscanf(s, "%1x%1x%1x", &c.R, &c.G, &c.B)
		c.R *= 17
		c.G *= 17
		c.B *= 17
	}
	return c
}

func colorToScale(clr color.Color) (float64, float64, float64, float64) {
	r, g, b, a := clr.RGBA()
	rf := float64(r) / 0xffff
	gf := float64(g) / 0xffff
	bf := float64(b) / 0xffff
	af := float64(a) / 0xffff
	// Convert to non-premultiplied alpha components.
	if 0 < af {
		rf /= af
		gf /= af
		bf /= af
	}
	return rf, gf, bf, af
}
