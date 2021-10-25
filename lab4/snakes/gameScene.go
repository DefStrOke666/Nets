package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/borodun/nsu-nets/lab4/snakes/state"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math"
	"math/rand"
	"net"
	"time"
)

const (
	multicastAddr = "239.192.0.4:9192"
)

type GameScene struct {
	background *ebiten.Image

	state        *proto.GameState
	stateChanged bool
	canJoin      bool

	playerName   string
	playerSnakes map[string]*proto.GameState_Snake

	columns, rows   int
	fieldBackground *ebiten.Image
	field           *Field
	busyCells       [][]bool

	buttonPics []*Picture
	exit       bool

	maxID int

	lastUpdate time.Time
}

func NewGameScene(config *proto.GameConfig) *GameScene {
	scene := &GameScene{}

	scene.state = new(proto.GameState)
	scene.state.Config = config

	scene.stateChanged = false
	scene.canJoin = true
	scene.exit = false

	scene.columns = int(*scene.state.Config.Width)
	scene.rows = int(*scene.state.Config.Height)

	scene.buttonPics = make([]*Picture, 2)
	scene.busyCells = make([][]bool, scene.columns)
	for i := range scene.busyCells {
		scene.busyCells[i] = make([]bool, scene.rows)
	}

	scene.addFood(5)

	scene.maxID = 0
	scene.state.Players = &proto.GamePlayers{}
	scene.playerSnakes = make(map[string]*proto.GameState_Snake)
	scene.addPlayer("borodun", proto.PlayerType_HUMAN, false)

	scene.lastUpdate = time.Now()
	scene.updateImages()

	go scene.sendAnnouncement()
	return scene
}

func (g *GameScene) addPlayer(name string, pType proto.PlayerType, view bool) {
	player := state.CreatePlayer(name)
	*player.Type = pType
	if view {
		*player.Role = proto.NodeRole_VIEWER
	}
	g.maxID++
	*player.Id = int32(g.maxID)

	g.state.Players.Players = append(g.state.Players.Players, player)
	if !view {
		head, chk := g.findFreeSquare()
		if !chk {
			println("Couldn't find place fr snake, turning player into VIEWER")
			*player.Role = proto.NodeRole_VIEWER
			return
		}

		snake := state.CreateSnake(g.maxID, head)
		g.busyCells[head.GetX()][head.GetY()] = true
		switch snake.GetHeadDirection() {
		case proto.Direction_UP:
			tail := state.CreateCoord(0, -1)
			snake.Points = append(snake.Points, tail)
			x := (int32(g.columns) + head.GetX() - 1) % int32(g.columns)
			g.busyCells[x][head.GetY()] = true
			break
		case proto.Direction_DOWN:
			tail := state.CreateCoord(0, 1)
			snake.Points = append(snake.Points, tail)
			x := (int32(g.columns) + head.GetX() - 1) % int32(g.columns)
			g.busyCells[x][head.GetY()] = true
			break
		case proto.Direction_LEFT:
			tail := state.CreateCoord(1, 0)
			snake.Points = append(snake.Points, tail)
			y := (int32(g.rows) + head.GetY() - 1) % int32(g.rows)
			g.busyCells[head.GetX()][y] = true
			break
		case proto.Direction_RIGHT:
			tail := state.CreateCoord(-1, 0)
			snake.Points = append(snake.Points, tail)
			y := (int32(g.rows) + head.GetY() - 1) % int32(g.rows)
			g.busyCells[head.GetX()][y] = true
			break
		}
		g.state.Snakes = append(g.state.Snakes, snake)
		g.playerSnakes[name] = snake
	}
}

func (g *GameScene) findFreeSquare() (*proto.GameState_Coord, bool) {
	x, y := 0, 0
	randy := rand.New(rand.NewSource(time.Now().Unix()))
	found := false
	for i := 0; i < 10 && !found; i++ {
		x = randy.Intn(g.columns)
		y = randy.Intn(g.rows)
		if g.busyCells[x][y] == false {
			for X := -2; X < 3 && !found; X++ {
				for Y := -2; Y < 3 && !found; Y++ {
					fieldX := (g.columns + x + X) % g.columns
					fieldY := (g.rows + y + Y) % g.rows
					if g.busyCells[fieldX][fieldY] == true {
						found = true
					}
				}
			}
		}
	}
	println("X:", x, "Y:", y)
	if found {
		return state.CreateCoord(x, y), true
	} else {
		return state.CreateCoord(-1, -1), false
	}
}

func (g *GameScene) receiveMessages() {

}

func (g *GameScene) addFood(count int) {
	x, y := 0, 0
	randy := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < count; i++ {
		foundEmpty := false
		for !foundEmpty {
			x = randy.Intn(g.columns)
			y = randy.Intn(g.rows)
			if g.busyCells[x][y] == false {
				foundEmpty = true
			}
		}
		g.state.Foods = append(g.state.Foods, state.CreateCoord(x, y))
		g.busyCells[x][y] = true
	}
	g.stateChanged = true
}

func (g *GameScene) sendAnnouncement() {
	annMsg := &proto.GameMessage_AnnouncementMsg{}
	annMsg.CanJoin = new(bool)

	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	for {
		annMsg.Config = g.state.Config
		annMsg.Players = g.state.Players
		*annMsg.CanJoin = g.canJoin
		marshal, err := annMsg.Marshal()
		if err != nil {
			log.Fatal(err)
		}
		_, err = c.Write(marshal)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(1 * time.Second)
		if g.exit {
			println("Stopped announcing")
			return
		}
	}
}

func (g *GameScene) updateImages() {
	margin := int(margin)
	spacingsV := margin * 3
	spacingsH := margin * 3

	widthUnit := (screenWidth - spacingsH) / 16
	heightUnit := (screenHeight - spacingsV) / 10

	fieldW := widthUnit * 12
	fieldH := heightUnit * 9

	buttonW := widthUnit * 3
	buttonH := heightUnit

	//scoreW := widthUnit * 4
	//scoreH := fieldH
	cellWidth := math.Min(float64(fieldW)/float64(g.columns), float64(fieldH)/float64(g.rows))

	g.field = NewField(fieldW, fieldH, g.columns, g.rows, cellWidth)
	g.fieldBackground = ebiten.NewImage(fieldW, fieldH)
	g.field.Draw(g.fieldBackground)

	g.background = getRoundRect(screenWidth, screenHeight, backgroundColor)

	g.buttonPics[0] = NewPicture(
		getRoundRectWithBorder(buttonW, buttonH, centreIdleColor, lineIdleColor),
		getRoundRectWithBorder(buttonW, buttonH, centreActiveColor, lineActiveColor))
	g.buttonPics[1] = NewPicture(
		getRoundRectWithBorder(buttonW, buttonH, centreIdleColor, lineIdleColor),
		getRoundRectWithBorder(buttonW, buttonH, centreActiveColor, lineActiveColor),
	).SetHandler(func(state *GameState) {
		g.exit = true
		state.SceneManager.GoTo(NewTitleScene())
	})
	g.buttonPics[0].SetRect(g.buttonPics[0].GetIdleImage().Bounds().Add(image.Point{X: margin, Y: fieldH + margin*2}))
	g.buttonPics[1].SetRect(g.buttonPics[1].GetIdleImage().Bounds().Add(image.Point{X: margin*2 + buttonW, Y: fieldH + margin*2}))
}

func (g *GameScene) changeSnakeDirection(snake *proto.GameState_Snake) {
	direction := snake.GetHeadDirection()
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if !(direction == proto.Direction_UP || direction == proto.Direction_DOWN) {
			*snake.HeadDirection = proto.Direction_UP
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if !(direction == proto.Direction_UP || direction == proto.Direction_DOWN) {
			*snake.HeadDirection = proto.Direction_DOWN
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if !(direction == proto.Direction_RIGHT || direction == proto.Direction_LEFT) {
			*snake.HeadDirection = proto.Direction_RIGHT
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if !(direction == proto.Direction_RIGHT || direction == proto.Direction_LEFT) {
			*snake.HeadDirection = proto.Direction_LEFT
		}
	}
}

func (g *GameScene) moveSnake(snake *proto.GameState_Snake) {
	head := snake.Points[0]
	if head == nil {
		return
	}
	lastX, lastY := int(snake.Points[0].GetX()), int(snake.Points[0].GetY())

	switch snake.GetHeadDirection() {
	case proto.Direction_UP:
		*head.X = int32(lastX)
		*head.Y = int32((g.rows + lastY + 1) % g.rows)
	case proto.Direction_DOWN:
		*head.X = int32(lastX)
		*head.Y = int32((g.rows + lastY - 1) % g.rows)
	case proto.Direction_RIGHT:
		*head.X = int32((g.columns + lastY + 1) % g.columns)
		*head.Y = int32(lastY)
	case proto.Direction_LEFT:
		*head.X = int32((g.columns + lastY - 1) % g.columns)
		*head.Y = int32(lastY)
	}

}

func (g *GameScene) Update(state *GameState) error {
	state.State = g.state
	if sizeChanged {
		g.updateImages()
	}

	for i := range g.buttonPics {
		if g.canJoin {
			g.buttonPics[i].Update(state)
		}
	}

	g.changeSnakeDirection(g.playerSnakes["borodun"])

	// Update game
	if time.Now().After(g.lastUpdate.Add(time.Millisecond * time.Duration(g.state.Config.GetStateDelayMs()))) {
		g.moveSnake(g.playerSnakes["borodun"])

		err := g.field.Update(state)
		if err != nil {
			return err
		}
		g.stateChanged = false

		g.lastUpdate = time.Now()
	}

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(fillColor)
	screen.DrawImage(g.background, &ebiten.DrawImageOptions{})

	g.field.Draw(g.fieldBackground)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(margin, margin)
	screen.DrawImage(g.fieldBackground, op)

	for i := range g.buttonPics {
		g.buttonPics[i].Draw(screen)
	}
}
