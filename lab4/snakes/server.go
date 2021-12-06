package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"log"
	"math/rand"
	"net"
	"time"
)

func (g *GameScene) sendAnnouncement() {
	annMsg := &proto.GameMessage_AnnouncementMsg{}
	annMsg.CanJoin = new(bool)

	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	println("Sending announcements to:", addr.String())
	for {
		annMsg.Config = g.state.Config
		annMsg.Players = g.state.Players
		*annMsg.CanJoin = g.canJoin
		marshal, err := annMsg.Marshal()
		if err != nil {
			log.Print(err)
		}
		_, err = c.Write(marshal)
		if err != nil {
			log.Print(err)
		}

		time.Sleep(1 * time.Second)
		if g.exit {
			println("Stopped announcing")
			return
		}
	}
}

func (g *GameScene) processMessages() {
	servAddr, err := net.ResolveUDPAddr("udp", "")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", servAddr)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.SetReadBuffer(maxDatagramSize)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, maxDatagramSize)
	println("Listening for messages on:", servAddr.String())
	for {
		time.Sleep(10 * time.Millisecond)
		read, addr, err := conn.ReadFromUDP(b)
		if err != nil {
			log.Print("ReadFromUDP failed:", err)
		}
		println("Connection from:", addr.String())

		message := &proto.GameMessage{}
		err = message.Unmarshal(b[:read])
		if err != nil {
			log.Print(err)
		}

		switch msg := message.Type.(type) {
		case *proto.GameMessage_Join:
			println("Join from addr:", addr.String())
			ack := msg.Join
			println("Ack:", ack.String())
			g.addPlayer(ack.GetName(), ack.GetPlayerType(), ack.GetOnlyView())
			answer := utils.CreateAckMessage(message.GetMsgSeq(), message.GetSenderId(), message.GetReceiverId())
			marshaledAnswer, err := answer.Marshal()
			if err != nil {
				log.Print(err)
			}
			_, err = conn.Write(marshaledAnswer)
			if err != nil {
				log.Print(err)
			}
			continue
		}

		if g.exit {
			println("Stopped listening for messages")
			return
		}
	}
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
