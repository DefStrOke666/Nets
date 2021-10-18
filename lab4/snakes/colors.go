package snakes

import (
	"fmt"
	"image/color"
)

var (
	fillColor = ParseHexColor("2A0944")

	backgroundColor = ParseHexColor("170055")

	centreIdleColor   = ParseHexColor("C9CCD5")
	centreActiveColor = ParseHexColor("FFE3E3")
	lineIdleColor     = ParseHexColor("FF5C58")
	lineActiveColor   = ParseHexColor("FDB827")

	titleIdleColor   = ParseHexColor("FFB319")
	titleActiveColor = ParseHexColor("E63E6D")
	idleColor        = ParseHexColor("80ED99")
	activeColor      = ParseHexColor("FF7777")

	serverBackgroundIdleColor   = ParseHexColor("082032")
	serverBackgroundActiveColor = ParseHexColor("334756")
	serverTextIdleColor         = ParseHexColor("FFF8E5")
	serverTextActiveColor       = ParseHexColor("FF4C29")

	fieldCellColor1 = ParseHexColor("FFD56B")
	fieldCellColor2 = ParseHexColor("FFB26B")
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
