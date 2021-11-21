package utils

import (
	"fmt"
	"image/color"
)

var (
	FillColor = ParseHexColor("2A0944")

	BackgroundColor = ParseHexColor("170055")

	CentreIdleColor   = ParseHexColor("C9CCD5")
	CentreActiveColor = ParseHexColor("FFE3E3")
	LineIdleColor     = ParseHexColor("FF5C58")
	LineActiveColor   = ParseHexColor("FDB827")

	TitleIdleColor   = ParseHexColor("FFB319")
	TitleActiveColor = ParseHexColor("E63E6D")
	IdleColor        = ParseHexColor("80ED99")
	ActiveColor      = ParseHexColor("FF7777")

	ServerBackgroundIdleColor   = ParseHexColor("082032")
	ServerBackgroundActiveColor = ParseHexColor("334756")
	ServerTextIdleColor         = ParseHexColor("FFF8E5")
	ServerTextActiveColor       = ParseHexColor("FF4C29")

	ScoreCentreColor = ParseHexColor("142F43")
	ScoreLineColor   = ParseHexColor("FFAB4C")
	ScoreTextColor   = ParseHexColor("99FEFF")

	ConfigCentreColor = ParseHexColor("142F43")
	ConfigLineColor   = ParseHexColor("FFAB4C")
	ConfigTextColor   = ParseHexColor("99FEFF")

	FieldCellColor1 = ParseHexColor("FFD56B")
	FieldCellColor2 = ParseHexColor("FFB26B")

	FoodColor       = ParseHexColor("E02401")
	SnakeBodyColor1 = ParseHexColor("FF8243")
	SnakeHeadColor1 = ParseHexColor("E26A2C")
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
