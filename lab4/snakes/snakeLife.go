package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"math/rand"
)

func (g *GameScene) getDir() (proto.Direction, bool) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		return proto.Direction_UP, true
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		return proto.Direction_DOWN, true
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		return proto.Direction_RIGHT, true
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		return proto.Direction_LEFT, true
	}
	return proto.Direction_UP, false
}

func (g *GameScene) changeSnakeDirection(name string, newDir proto.Direction) {
	snake := g.playerSnakes[name]
	oldDir := g.playerSaveDir[name]
	if snake == nil {
		return
	}
	switch newDir {
	case proto.Direction_UP:
		if !(oldDir == proto.Direction_UP || oldDir == proto.Direction_DOWN) {
			*snake.HeadDirection = proto.Direction_UP
		}
	case proto.Direction_DOWN:
		if !(oldDir == proto.Direction_UP || oldDir == proto.Direction_DOWN) {
			*snake.HeadDirection = proto.Direction_DOWN
		}
	case proto.Direction_RIGHT:
		if !(oldDir == proto.Direction_RIGHT || oldDir == proto.Direction_LEFT) {
			*snake.HeadDirection = proto.Direction_RIGHT
		}
	case proto.Direction_LEFT:
		if !(oldDir == proto.Direction_RIGHT || oldDir == proto.Direction_LEFT) {
			*snake.HeadDirection = proto.Direction_LEFT
		}
	}
}

func (g *GameScene) moveSnake(snake *proto.GameState_Snake) {
	if snake == nil {
		return
	}
	head := snake.Points[0]
	if head == nil {
		return
	}

	prevHeadX, prevHeadY := int(head.GetX()), int(head.GetY())

	// move head
	var newPoint *proto.GameState_Coord
	switch snake.GetHeadDirection() {
	case proto.Direction_UP:
		*head.X = int32(prevHeadX)
		*head.Y = int32((g.rows + prevHeadY - 1) % g.rows)
		newPoint = utils.CreateCoord(0, 1)
		break
	case proto.Direction_DOWN:
		*head.X = int32(prevHeadX)
		*head.Y = int32((g.rows + prevHeadY + 1) % g.rows)
		newPoint = utils.CreateCoord(0, -1)
		break
	case proto.Direction_RIGHT:
		*head.X = int32((g.columns + prevHeadX + 1) % g.columns)
		*head.Y = int32(prevHeadY)
		newPoint = utils.CreateCoord(-1, 0)
		break
	case proto.Direction_LEFT:
		*head.X = int32((g.columns + prevHeadX - 1) % g.columns)
		*head.Y = int32(prevHeadY)
		newPoint = utils.CreateCoord(1, 0)
		break
	}

	newHeadX, newHeadY := int(head.GetX()), int(head.GetY())

	//move neck
	neck := snake.Points[1]
	pointX, pointY := int(neck.GetX()), int(neck.GetY())
	dX, dY := prevHeadX+pointX, prevHeadY+pointY
	pointsLen := len(snake.Points)
	if dX != newHeadX && dY != newHeadY {
		if math.Abs(float64(pointX+pointY)) > 1 {
			snake.Points = append(snake.Points[:2], snake.Points[1:]...)
			snake.Points[1] = newPoint
			pointsLen = len(snake.Points)
		} else {
			snake.Points[1] = newPoint
		}
	} else if pointsLen > 2 {
		if pointX != 0 {
			*neck.X = MoveFromZero(pointX, 1)
		} else {
			*neck.Y = MoveFromZero(pointY, 1)
		}
	}

	tail := snake.Points[pointsLen-1]
	tailX, tailY := int(tail.GetX()), int(tail.GetY())

	// move tail
	if pointsLen > 2 {
		if math.Abs(float64(tailX+tailY)) > 1 {
			if tailX != 0 {
				*tail.X = MoveToZero(tailX, 1)
			} else {
				*tail.Y = MoveToZero(tailY, 1)
			}
		} else {
			snake.Points = snake.Points[:pointsLen-1]
		}
	}
}

func (g *GameScene) clearSnakeCells() {
	for i, row := range g.snakeCells {
		for j := range row {
			g.snakeCells[i][j] = false
		}
	}
}

func (g *GameScene) fillSnakeCells(snake *proto.GameState_Snake) {
	if snake == nil {
		return
	}
	head := snake.Points[0]
	if head == nil {
		return
	}

	lastX, lastY := head.GetX(), head.GetY()
	for i, point := range snake.Points {
		if i == 0 {
			g.snakeCells[lastX][lastY] = false
		} else {
			pointX, pointY := point.GetX(), point.GetY()
			lineLength := int(math.Abs(float64(pointX + pointY)))
			lastX, lastY = lastX+pointX, lastY+pointY

			for j := 0; j < lineLength; j++ {
				if pointX != 0 {
					pointX = MoveToZero(int(pointX), 1)
				} else {
					pointY = MoveToZero(int(pointY), 1)
				}
				xCoord := (int32(g.columns) + lastX - pointX) % int32(g.columns)
				yCoord := (int32(g.rows) + lastY - pointY) % int32(g.rows)
				g.snakeCells[xCoord][yCoord] = true
			}
		}
	}
}

func (g *GameScene) makeFoodFromSnake(snake *proto.GameState_Snake) {
	if snake == nil {
		return
	}
	head := snake.Points[0]
	if head == nil {
		return
	}

	lastX, lastY := snake.Points[0].GetX(), snake.Points[0].GetY()
	for i, point := range snake.Points {
		if i == 0 {
			g.foodCells[lastX][lastY] = true
			g.state.Foods = append(g.state.Foods, utils.CreateCoord(int(lastX), int(lastY)))
		} else {
			pointX, pointY := point.GetX(), point.GetY()
			lineLength := int(math.Abs(float64(pointX + pointY)))
			lastX, lastY = lastX+pointX, lastY+pointY

			for j := 0; j < lineLength; j++ {
				randInt := rand.Intn(100)
				threshold := g.state.Config.GetDeadFoodProb()
				if float32(randInt)/100 <= threshold {
					if pointX != 0 {
						pointX = MoveToZero(int(pointX), 1)
					} else {
						pointY = MoveToZero(int(pointY), 1)
					}
					xCoord := (int32(g.columns) + lastX - pointX) % int32(g.columns)
					yCoord := (int32(g.rows) + lastY - pointY) % int32(g.rows)
					g.foodCells[xCoord][yCoord] = true
					g.state.Foods = append(g.state.Foods, utils.CreateCoord(int(xCoord), int(yCoord)))
				}
			}
		}
	}
}

func (g *GameScene) checkCollision(snake *proto.GameState_Snake) bool {
	if snake == nil {
		return false
	}
	head := snake.Points[0]
	if head == nil {
		return false
	}

	headX, headY := int(head.GetX()), int(head.GetY())
	return g.snakeCells[headX][headY]
}

func MoveFromZero(x, step int) int32 {
	if x > 0 {
		return int32(x + step)
	}
	return int32(x - step)
}

func MoveToZero(x, step int) int32 {
	if x > 0 {
		return int32(x - step)
	}
	return int32(x + step)
}

func (g *GameScene) eatFood(snake *proto.GameState_Snake, name string) {
	if snake == nil {
		return
	}
	head := snake.Points[0]
	if head == nil {
		return
	}

	headX, headY := int(head.GetX()), int(head.GetY())
	if g.foodCells[headX][headY] == true {
		pointsLen := len(snake.Points)
		tail := snake.Points[pointsLen-1]
		tailX, tailY := int(tail.GetX()), int(tail.GetY())

		if tailX != 0 {
			*tail.X = MoveFromZero(tailX, 1)
		} else {
			*tail.Y = MoveFromZero(tailY, 1)
		}

		g.foodCells[headX][headY] = false
		for i, coord := range g.state.Foods {
			if int(coord.GetX()) == headX && int(coord.GetY()) == headY {
				g.state.Foods = append(g.state.Foods[:i], g.state.Foods[i+1:]...)
				break
			}
		}
		g.ateFood = true
		g.addFood(1)
		*g.playersByName[name].Score++
	}
}

func (g *GameScene) removeSnake(snake *proto.GameState_Snake, name string) {
	index := -1
	for i, snake2 := range g.state.Snakes {
		if snake2 == snake {
			index = i
		}
	}
	*snake.State = proto.GameState_Snake_ZOMBIE
	g.state.Snakes = append(g.state.Snakes[:index], g.state.Snakes[index+1:]...)
	delete(g.playerSnakes, name)
}
