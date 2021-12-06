package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math"
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

	isServer     bool
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

	acks []*proto.GameMessage_Ack
	conn net.Conn

	buttonPics []*utils.Picture
	exit       bool

	maxID int

	lastUpdate time.Time
}

func NewGameScene(config *proto.GameConfig, serverAddr *net.UDPAddr, view bool) *GameScene {
	scene := &GameScene{}

	scene.state = new(proto.GameState)
	scene.state.Config = config
	scene.columns = int(*scene.state.Config.Width)
	scene.rows = int(*scene.state.Config.Height)
	scene.buttonPics = make([]*utils.Picture, 2)
	scene.snakeCells = make([][]bool, scene.columns)

	for i := range scene.snakeCells {
		scene.snakeCells[i] = make([]bool, scene.rows)
	}

	scene.foodCells = make([][]bool, scene.columns)
	for i := range scene.foodCells {
		scene.foodCells[i] = make([]bool, scene.rows)
	}

	scene.state.Players = &proto.GamePlayers{}
	scene.playersByName = make(map[string]*proto.GamePlayer)
	scene.playerSnakes = make(map[string]*proto.GameState_Snake)
	scene.playerSaveDir = make(map[string]proto.Direction)

	scene.lastUpdate = time.Now()
	scene.updateImages()

	if serverAddr != nil {
		scene.playerName = utils.Conf.PlayerNames.PlayerName
		scene.isServer = false
		conn := scene.joinServer(serverAddr, view)
		if conn == nil {
			log.Fatalln("Connection error")
		} else {
			scene.conn = conn
		}
	} else {
		scene.playerName = utils.Conf.PlayerNames.AdminName
		scene.isServer = true
		scene.stateChanged = false
		scene.canJoin = true
		scene.exit = false
		scene.addFood(int(config.GetFoodStatic()))
		scene.maxID = 0
		scene.addPlayer(scene.playerName, proto.PlayerType_HUMAN, false)
		go scene.sendAnnouncement()
		go scene.processMessages()
	}

	return scene
}

func (g *GameScene) updateImages() {
	margin := int(utils.Margin)
	spacingsV := margin*3 + int(utils.LineThickness*2)
	spacingsH := margin*3 + int(utils.LineThickness*2)

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
	g.scoreH = actialH + int(utils.LineThickness*2)
	g.drawScore()
	g.drawConfig()

	g.field = NewField(g.columns, g.rows, cellWidth)
	g.fieldBackground = utils.GetRectWithBorder(actialW+int(utils.LineThickness*2), actialH+int(utils.LineThickness*2), utils.CentreActiveColor, utils.LineActiveColor)
	g.field.Draw(g.fieldBackground)

	g.background = utils.GetRoundRect(screenWidth, screenHeight, utils.BackgroundColor)

	g.buttonPics[0] = utils.NewPicture(
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreIdleColor, utils.LineIdleColor, "View", utils.GetMenuFonts(4)),
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreActiveColor, utils.LineActiveColor, "View", utils.GetMenuFonts(4)))
	g.buttonPics[1] = utils.NewPicture(
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreIdleColor, utils.LineIdleColor, "Exit", utils.GetMenuFonts(4)),
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreActiveColor, utils.LineActiveColor, "Exit", utils.GetMenuFonts(4)),
	).SetHandler(func() {
		g.exit = true
		sceneManager.GoTo(NewTitleScene())
	})
	g.buttonPics[0].SetRect(g.buttonPics[0].GetIdleImage().Bounds().Add(image.Point{X: margin, Y: fieldH + margin*2 + int(utils.LineThickness*2)}))
	g.buttonPics[1].SetRect(g.buttonPics[1].GetIdleImage().Bounds().Add(image.Point{X: margin*2 + buttonW, Y: fieldH + margin*2 + int(utils.LineThickness*2)}))
}

func (g *GameScene) drawScore() {
	namesImg := ebiten.NewImage(utils.TextWidth("VeryLongName", utils.GetMenuFonts(3)), g.scoreH)
	numsImg := ebiten.NewImage(utils.TextWidth("9999", utils.GetMenuFonts(3)), g.scoreH)

	op := &ebiten.DrawImageOptions{}
	bckImg := utils.GetRoundRectWithBorder(g.scoreW, g.scoreH, utils.ScoreCentreColor, utils.ScoreLineColor)
	for _, player := range g.state.Players.GetPlayers() {
		score := strconv.Itoa(int(player.GetScore()))
		name := player.GetName()
		textH := utils.TextHeight(name+score, utils.GetMenuFonts(3))
		namesImg.DrawImage(utils.CreateStringImage(name, utils.GetMenuFonts(3), utils.ScoreTextColor), op)
		numsImg.DrawImage(utils.CreateStringImage(score, utils.GetMenuFonts(3), utils.ScoreTextColor), op)
		op.GeoM.Translate(0, float64(textH)+utils.Margin)
	}
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(utils.Margin, utils.Margin)
	bckImg.DrawImage(namesImg, op2)
	op2.GeoM.Translate(float64(namesImg.Bounds().Max.X), 0)
	bckImg.DrawImage(numsImg, op2)
	g.scoreImg = bckImg
}

func (g *GameScene) drawConfig() {
	img := utils.GetRectWithBorder(g.scoreW, g.scoreH, utils.ConfigCentreColor, utils.ConfigLineColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(utils.Margin, utils.Margin)

	configStr := strings.Split(g.state.Config.String(), ",")
	for _, s := range configStr {
		if s != "" {
			textH := utils.TextHeight(s, utils.GetMenuFonts(3))
			img.DrawImage(utils.CreateStringImage(s, utils.GetMenuFonts(3), utils.ConfigTextColor), op)
			op.GeoM.Translate(0, float64(textH)+utils.Margin)
		}
	}

	g.configImg = img
}

func (g *GameScene) Update(state *GameState) error {
	state.State = g.state
	if sizeChanged {
		g.updateImages()
	}

	if g.isServer {
		for i := range g.buttonPics {
			g.buttonPics[i].Update()
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
						g.removeSnake(snake, name)
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
	}

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(utils.FillColor)
	screen.DrawImage(g.background, &ebiten.DrawImageOptions{})

	g.field.Draw(g.fieldBackground)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(utils.Margin, utils.Margin)
	screen.DrawImage(g.fieldBackground, op)

	for i := range g.buttonPics {
		g.buttonPics[i].Draw(screen)
	}

	op.GeoM.Reset()
	op.GeoM.Translate(float64(screenWidth-g.scoreW-int(utils.Margin)), utils.Margin)
	screen.DrawImage(g.scoreImg, op)
	op.GeoM.Translate(float64(-g.scoreW-int(utils.Margin)), 0)
	screen.DrawImage(g.configImg, op)
}
