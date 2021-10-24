package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type Field struct {
	width, height int
	columns, rows int
	cellWidth     float64

	emptyField *ebiten.Image
	field      *ebiten.Image
	foodImg    *ebiten.Image
}

func NewField(w, h, c, r int, cw float64) *Field {
	field := &Field{}
	field.width = w
	field.height = h
	field.columns = c
	field.rows = r
	field.cellWidth = cw
	field.emptyField = ebiten.NewImage(w, h)

	field.clearField()
	field.createFood()

	op := &ebiten.DrawImageOptions{}
	field.emptyField.DrawImage(field.field, op)
	return field
}

func (f *Field) clearField() {
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

func (f *Field) createFood() {
	dc := gg.NewContext(int(f.cellWidth), int(f.cellWidth))
	dc.DrawRectangle(0, 0, f.cellWidth, f.cellWidth)
	dc.SetColor(foodColor)
	dc.Fill()
	f.foodImg = ebiten.NewImageFromImage(dc.Image())
}

func (f *Field) drawFood(food []*proto.GameState_Coord) {
	op := &ebiten.DrawImageOptions{}
	for _, coord := range food {
		op.GeoM.Reset()
		op.GeoM.Translate(float64(coord.GetX())*f.cellWidth, float64(coord.GetY())*f.cellWidth)
		f.field.DrawImage(f.foodImg, op)
	}
}

func (f *Field) drawSnakes(snakes []*proto.GameState_Snake) {
	println("DrawSnakes:", snakes)
	for _, snake := range snakes {
		f.drawSnake(snake)
	}
}

func (f *Field) drawSnake(snake *proto.GameState_Snake) {
	println("DrawSnake:", snake)
	dc := gg.NewContext(int(f.cellWidth), int(f.cellWidth))
	dc.DrawRectangle(0, 0, f.cellWidth, f.cellWidth)
	dc.SetColor(snakeBodyColor1)
	dc.Fill()
	snakeCell := ebiten.NewImageFromImage(dc.Image())

	dc = gg.NewContext(int(f.cellWidth), int(f.cellWidth))
	dc.DrawRectangle(0, 0, f.cellWidth, f.cellWidth)
	dc.SetColor(snakeHeadColor1)
	dc.Fill()
	snakeHead := ebiten.NewImageFromImage(dc.Image())

	op := &ebiten.DrawImageOptions{}

	lastX, lastY := float64(snake.Points[0].GetX()), float64(snake.Points[0].GetY())
	for i, point := range snake.Points {
		op.GeoM.Reset()
		if i == 0 {
			op.GeoM.Translate(lastX*f.cellWidth, lastY*f.cellWidth)
			f.field.DrawImage(snakeHead, op)
		} else {
			op.GeoM.Translate((lastX+float64(point.GetX()))*f.cellWidth, (lastY+float64(point.GetY()))*f.cellWidth)
			f.field.DrawImage(snakeCell, op)
			lastX, lastY = lastX+float64(point.GetX()), lastY+float64(point.GetY())
		}
	}
}

func (f *Field) Update(state *GameState) error {
	op := &ebiten.DrawImageOptions{}
	f.field.DrawImage(f.emptyField, op)

	f.drawFood(state.State.Foods)
	f.drawSnakes(state.State.Snakes)
	return nil
}

func (f *Field) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(f.field, op)
}
