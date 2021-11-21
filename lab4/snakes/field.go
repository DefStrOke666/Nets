package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Field struct {
	width, height int
	columns, rows int
	cellWidth     int

	emptyField *ebiten.Image
	field      *ebiten.Image
	foodImg    *ebiten.Image
}

func NewField(c, r, cw int) *Field {
	field := &Field{}
	field.columns = c
	field.rows = r
	field.cellWidth = cw
	field.width = c * cw
	field.height = r * cw
	field.emptyField = ebiten.NewImage(field.width, field.height)

	field.createField()
	field.createFood()

	op := &ebiten.DrawImageOptions{}
	field.emptyField.DrawImage(field.field, op)
	return field
}

func (f *Field) createField() {
	dc := gg.NewContext(f.width, f.height)
	for y := 0; y < f.rows; y++ {
		for x := 0; x < f.columns; x++ {
			dc.DrawRectangle(float64(x*f.cellWidth), float64(y*f.cellWidth), float64(f.cellWidth), float64(f.cellWidth))
			if (x+y)%2 == 0 {
				dc.SetColor(utils.FieldCellColor1)
			} else {
				dc.SetColor(utils.FieldCellColor2)
			}
			dc.Fill()
		}
	}
	f.field = ebiten.NewImageFromImage(dc.Image())
}

func (f *Field) createFood() {
	dc := gg.NewContext(int(f.cellWidth), int(f.cellWidth))
	dc.DrawRectangle(0, 0, float64(f.cellWidth), float64(f.cellWidth))
	dc.SetColor(utils.FoodColor)
	dc.Fill()
	f.foodImg = ebiten.NewImageFromImage(dc.Image())
}

func (f *Field) drawFood(food []*proto.GameState_Coord) {
	op := &ebiten.DrawImageOptions{}
	for _, coord := range food {
		op.GeoM.Reset()
		op.GeoM.Translate(float64(int(coord.GetX())*f.cellWidth), float64(int(coord.GetY())*f.cellWidth))
		f.field.DrawImage(f.foodImg, op)
	}
}

func (f *Field) drawSnakes(snakes []*proto.GameState_Snake) {
	for _, snake := range snakes {
		f.drawSnake(snake)
	}
}

func (f *Field) drawSnake(snake *proto.GameState_Snake) {
	dc := gg.NewContext(int(f.cellWidth), int(f.cellWidth))
	dc.DrawRectangle(0, 0, float64(f.cellWidth), float64(f.cellWidth))
	dc.SetColor(utils.SnakeBodyColor1)
	dc.Fill()
	snakeCell := ebiten.NewImageFromImage(dc.Image())

	dc = gg.NewContext(int(f.cellWidth), int(f.cellWidth))
	dc.DrawRectangle(0, 0, float64(f.cellWidth), float64(f.cellWidth))
	dc.SetColor(utils.SnakeHeadColor1)
	dc.Fill()
	snakeHead := ebiten.NewImageFromImage(dc.Image())

	op := &ebiten.DrawImageOptions{}

	lastX, lastY := snake.Points[0].GetX(), snake.Points[0].GetY()
	for i, point := range snake.Points {
		op.GeoM.Reset()
		if i == 0 {
			op.GeoM.Translate(float64(int(lastX)*f.cellWidth), float64(int(lastY)*f.cellWidth))
			f.field.DrawImage(snakeHead, op)
		} else {
			pointX, pointY := point.GetX(), point.GetY()
			lineLength := int(math.Abs(float64(pointX + pointY)))
			lastX, lastY = lastX+pointX, lastY+pointY

			for j := 0; j < lineLength; j++ {
				op.GeoM.Reset()
				if pointX != 0 {
					pointX = MoveToZero(int(pointX), 1)
				} else {
					pointY = MoveToZero(int(pointY), 1)
				}
				xCoord := (int32(f.columns) + lastX - pointX) % int32(f.columns)
				yCoord := (int32(f.rows) + lastY - pointY) % int32(f.rows)
				op.GeoM.Translate(float64(int(xCoord)*f.cellWidth), float64(int(yCoord)*f.cellWidth))
				f.field.DrawImage(snakeCell, op)
			}
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
	op.GeoM.Translate(utils.LineThickness, utils.LineThickness)
	screen.DrawImage(f.field, op)
}
