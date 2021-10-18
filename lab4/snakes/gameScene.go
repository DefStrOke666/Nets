package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/borodun/nsu-nets/lab4/snakes/state"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math"
	"net"
	"time"
)

const (
	multicastAddr = "239.192.0.4:9192"
)

type GameScene struct {
	background *ebiten.Image

	state   *proto.GameState
	canJoin bool

	columns, rows   int
	fieldBackground *ebiten.Image
	field           *Field
}

func NewGameScene(config *proto.GameConfig) *GameScene {
	scene := &GameScene{}

	scene.state = new(proto.GameState)
	scene.state.Config = config
	players := &proto.GamePlayers{}
	for i := 0; i < 5; i++ {
		players.Players = append(players.Players, state.CreatePlayer("player"))
	}
	scene.state.Players = players
	scene.canJoin = true

	scene.columns = int(*scene.state.Config.Width)
	scene.rows = int(*scene.state.Config.Height)

	go scene.sendAnnouncement()

	scene.updateImages()
	return scene
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

	//scoreW := widthUnit * 4
	//scoreH := fieldH
	cellWidth := math.Min(float64(fieldW)/float64(g.columns), float64(fieldH)/float64(g.rows))

	g.field = NewField(fieldW, fieldH, g.columns, g.rows, cellWidth)
	g.fieldBackground = ebiten.NewImage(fieldW, fieldH)
	g.field.Draw(g.fieldBackground)

	g.background = getRoundRect(screenWidth, screenHeight, backgroundColor)
}

func (g *GameScene) Update(state *GameState) error {
	if sizeChanged {
		g.updateImages()
	}

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(fillColor)
	screen.DrawImage(g.background, &ebiten.DrawImageOptions{})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(margin, margin)
	screen.DrawImage(g.fieldBackground, op)
}
