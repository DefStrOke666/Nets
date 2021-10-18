package snakes

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type Field struct {
	width, height int
	columns, rows int
	cellWidth     float64

	field *ebiten.Image
}

func NewField(w, h, c, r int, cw float64) *Field {
	field := &Field{}
	field.width = w
	field.height = h
	field.columns = c
	field.rows = r
	field.cellWidth = cw

	field.createField()
	return field
}

func (f *Field) createField() {
	dc := gg.NewContext(f.width, f.height)
	for y := 0; y < f.rows; y++ {
		for x := 0; x < f.columns; x++ {
			dc.DrawRectangle(float64(x)*f.cellWidth, float64(y)*f.cellWidth, f.cellWidth, f.cellWidth)
			if (x+y)%2 == 0 {
				dc.SetColor(fieldCellColor1)
			} else {
				dc.SetColor(fieldCellColor2)
			}
			dc.Fill()
		}
	}
	f.field = ebiten.NewImageFromImage(dc.Image())
}

func (f *Field) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(f.field, op)
}
