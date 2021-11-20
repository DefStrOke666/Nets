package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
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

	playerName    string
	playersByName map[string]*proto.GamePlayer
	playerSnakes  map[string]*proto.GameState_Snake
	playerSaveDir map[string]proto.Direction

	columns, rows   int
	fieldBackground *ebiten.Image
	field           *Field
	snakeCells      [][]bool
	foodCells       [][]bool
	ateFood         bool

	configImg *ebiten.Image
	scoreImg  *ebiten.Image
	scoreW    int
	scoreH    int

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
	scene.snakeCells = make([][]bool, scene.columns)
	for i := range scene.snakeCells {
		scene.snakeCells[i] = make([]bool, scene.rows)
	}

	scene.foodCells = make([][]bool, scene.columns)
	for i := range scene.foodCells {
		scene.foodCells[i] = make([]bool, scene.rows)
	}

	scene.addFood(int(config.GetFoodStatic()))

	scene.maxID = 0
	scene.state.Players = &proto.GamePlayers{}
	scene.playersByName = make(map[string]*proto.GamePlayer)
	scene.playerSnakes = make(map[string]*proto.GameState_Snake)
	scene.playerSaveDir = make(map[string]proto.Direction)
	scene.addPlayer("borodun", proto.PlayerType_HUMAN, false)

	scene.lastUpdate = time.Now()
	scene.updateImages()

	go scene.sendAnnouncement()
	return scene
}

func (g *GameScene) updateImages() {
	margin := int(margin)
	spacingsV := margin*3 + int(lineThickness*2)
	spacingsH := margin*3 + int(lineThickness*2)

	widthUnit := (screenWidth - spacingsH) / 16
	heightUnit := (screenHeight - spacingsV) / 10

	fieldW := widthUnit * 10
	fieldH := heightUnit * 9

	cellWidth := int(math.Min((float64(fieldW))/float64(g.columns), (float64(fieldH))/float64(g.rows)))
	actialW := g.columns * cellWidth
	actialH := g.rows * cellWidth

	buttonW := actialW / 2
	buttonH := heightUnit

	g.scoreW = (screenWidth - spacingsH - actialW - margin) / 2
	g.scoreH = actialH + int(lineThickness*2)
	g.drawScore()
	g.drawConfig()

	g.field = NewField(g.columns, g.rows, cellWidth)
	g.fieldBackground = getRectWithBorder(actialW+int(lineThickness*2), actialH+int(lineThickness*2), centreActiveColor, lineActiveColor)
	g.field.Draw(g.fieldBackground)

	g.background = getRoundRect(screenWidth, screenHeight, backgroundColor)

	g.buttonPics[0] = NewPicture(
		borderedRoundRectWithText(buttonW, buttonH, centreIdleColor, lineIdleColor, "View", getMenuFonts(4)),
		borderedRoundRectWithText(buttonW, buttonH, centreActiveColor, lineActiveColor, "View", getMenuFonts(4)))
	g.buttonPics[1] = NewPicture(
		borderedRoundRectWithText(buttonW, buttonH, centreIdleColor, lineIdleColor, "Exit", getMenuFonts(4)),
		borderedRoundRectWithText(buttonW, buttonH, centreActiveColor, lineActiveColor, "Exit", getMenuFonts(4)),
	).SetHandler(func(state *GameState) {
		g.exit = true
		state.SceneManager.GoTo(NewTitleScene())
	})
	g.buttonPics[0].SetRect(g.buttonPics[0].GetIdleImage().Bounds().Add(image.Point{X: margin, Y: fieldH + margin*2 + int(lineThickness*2)}))
	g.buttonPics[1].SetRect(g.buttonPics[1].GetIdleImage().Bounds().Add(image.Point{X: margin*2 + buttonW, Y: fieldH + margin*2 + int(lineThickness*2)}))
}

func (g *GameScene) drawScore() {
	namesImg := ebiten.NewImage(textWidth("VeryLongName", getMenuFonts(3)), g.scoreH)
	numsImg := ebiten.NewImage(textWidth("9999", getMenuFonts(3)), g.scoreH)

	op := &ebiten.DrawImageOptions{}
	bckImg := getRoundRectWithBorder(g.scoreW, g.scoreH, scoreCentreColor, scoreLineColor)
	for _, player := range g.state.Players.GetPlayers() {
		score := strconv.Itoa(int(player.GetScore()))
		name := player.GetName()
		textH := textHeight(name+score, getMenuFonts(3))
		namesImg.DrawImage(createStringImage(name, getMenuFonts(3), scoreTextColor), op)
		numsImg.DrawImage(createStringImage(score, getMenuFonts(3), scoreTextColor), op)
		op.GeoM.Translate(0, float64(textH)+margin)
	}
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(margin, margin)
	bckImg.DrawImage(namesImg, op2)
	op2.GeoM.Translate(float64(namesImg.Bounds().Max.X), 0)
	bckImg.DrawImage(numsImg, op2)
	g.scoreImg = bckImg
}

func (g *GameScene) drawConfig() {
	img := getRectWithBorder(g.scoreW, g.scoreH, configCentreColor, configLineColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(margin, margin)

	configStr := strings.Split(g.state.Config.String(), ",")
	for _, s := range configStr {
		if s != "" {
			textH := textHeight(s, getMenuFonts(3))
			img.DrawImage(createStringImage(s, getMenuFonts(3), configTextColor), op)
			op.GeoM.Translate(0, float64(textH)+margin)
		}
	}

	g.configImg = img
}

func (g *GameScene) addPlayer(name string, pType proto.PlayerType, view bool) {
	player := utils.CreatePlayer(name)
	*player.Type = pType
	if view {
		*player.Role = proto.NodeRole_VIEWER
	}
	g.maxID++
	*player.Id = int32(g.maxID)

	g.state.Players.Players = append(g.state.Players.Players, player)
	g.playersByName[name] = player
	if !view {
		head, chk := g.findFreeSquare()
		if !chk {
			println("Couldn't find place for snake, turning player into VIEWER")
			*player.Role = proto.NodeRole_VIEWER
			return
		}

		snake := utils.CreateSnake(g.maxID, head)
		g.snakeCells[head.GetX()][head.GetY()] = true
		var tail *proto.GameState_Coord
		switch snake.GetHeadDirection() {
		case proto.Direction_UP:
			tail = utils.CreateCoord(0, 1)
			snake.Points = append(snake.Points, tail)
			x := (int32(g.columns) + head.GetX() - 1) % int32(g.columns)
			g.snakeCells[x][head.GetY()] = true
			break
		case proto.Direction_DOWN:
			tail := utils.CreateCoord(0, -1)
			snake.Points = append(snake.Points, tail)
			x := (int32(g.columns) + head.GetX() - 1) % int32(g.columns)
			g.snakeCells[x][head.GetY()] = true
			break
		case proto.Direction_LEFT:
			tail := utils.CreateCoord(1, 0)
			snake.Points = append(snake.Points, tail)
			y := (int32(g.rows) + head.GetY() - 1) % int32(g.rows)
			g.snakeCells[head.GetX()][y] = true
			break
		case proto.Direction_RIGHT:
			tail := utils.CreateCoord(-1, 0)
			snake.Points = append(snake.Points, tail)
			y := (int32(g.rows) + head.GetY() - 1) % int32(g.rows)
			g.snakeCells[head.GetX()][y] = true
			break
		}
		g.state.Snakes = append(g.state.Snakes, snake)
		g.playerSnakes[name] = snake
		g.playerSaveDir[name] = *snake.HeadDirection
		g.addFood(int(g.state.Config.GetFoodPerPlayer()))
	}
}

func (g *GameScene) findFreeSquare() (*proto.GameState_Coord, bool) {
	x, y := 0, 0
	randy := rand.New(rand.NewSource(time.Now().Unix()))
	found := false
	for i := 0; i < 10 && !found; i++ {
		x = randy.Intn(g.columns)
		y = randy.Intn(g.rows)
		if g.snakeCells[x][y] == false && g.foodCells[x][y] == false {
			found = true
			for X := -2; X < 3; X++ {
				for Y := -2; Y < 3; Y++ {
					fieldX := (g.columns + x + X) % g.columns
					fieldY := (g.rows + y + Y) % g.rows
					if g.snakeCells[fieldX][fieldY] == true {
						found = false
					}
				}
			}
		}
	}
	println("X:", x, "Y:", y)
	if found {
		return utils.CreateCoord(x, y), true
	} else {
		return utils.CreateCoord(-1, -1), false
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
			if g.snakeCells[x][y] == false && g.foodCells[x][y] == false {
				foundEmpty = true
			}
		}
		g.state.Foods = append(g.state.Foods, utils.CreateCoord(x, y))
		g.foodCells[x][y] = true
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

	for name, snake := range g.playerSnakes {
		g.changeSnakeDirection(snake, name)
	}

	// Update game
	if time.Now().After(g.lastUpdate.Add(time.Millisecond * time.Duration(g.state.Config.GetStateDelayMs()))) {
		g.clearSnakeCells()
		for name, snake := range g.playerSnakes {
			if snake.GetState() == proto.GameState_Snake_ALIVE {
				g.moveSnake(snake)
				g.eatFood(snake, name)
				if g.ateFood {
					g.drawScore()
				}
				g.fillSnakeCells(snake)
				if g.checkCollision(snake) {
					println("Removing snake")
					g.makeFoodFromSnake(snake)
					g.removeSnake(snake)
				}
				g.playerSaveDir[name] = *snake.HeadDirection
			}
		}

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

	op.GeoM.Reset()
	op.GeoM.Translate(float64(screenWidth-g.scoreW-int(margin)), margin)
	screen.DrawImage(g.scoreImg, op)
	op.GeoM.Translate(float64(-g.scoreW-int(margin)), 0)
	screen.DrawImage(g.configImg, op)
}
