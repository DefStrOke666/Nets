package snakes

import (
	"fmt"
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
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
	servAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+g.port)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", servAddr)
	if err != nil {
		log.Fatal("ListenUDP:", err)
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
			continue
		}
		println("Connection from:", addr.String())

		message := &proto.GameMessage{}
		err = message.Unmarshal(b[:read])
		if err != nil {
			log.Print("Unmarshall error:", err)
			continue
		}

		switch msg := message.Type.(type) {
		case *proto.GameMessage_Join:
			println("Join from addr:", addr.String())
			join := msg.Join
			println("Join:", join.String())
			clientAddr := strings.Split(addr.String(), ":")
			ip := clientAddr[0]
			port, _ := strconv.Atoi(clientAddr[1])
			println("Client port:", port)
			g.addPlayer(join.GetName(), join.GetPlayerType(), join.GetOnlyView(), ip, port)
			g.namesById[int(message.GetSenderId())] = join.GetName()

			answer := utils.CreateGameMessage(1, 1, 2)
			answer.Type = &proto.GameMessage_Ack{Ack: &proto.GameMessage_AckMsg{}}
			marshaledAnswer, err := answer.Marshal()
			if err != nil {
				log.Print("marshall error:", err)
				continue
			}
			_, err = conn.WriteTo(marshaledAnswer, addr)
			if err != nil {
				log.Print("Write error:", err)
				continue
			}
			continue
		case *proto.GameMessage_Steer:
			println("Steer from addr:", addr.String())
			steer := msg.Steer
			println("Steer:", steer.String())
			name := g.namesById[int(message.GetSenderId())]
			g.changeSnakeDirection(name, steer.GetDirection())
		}

		if g.exit {
			println("Stopped listening for messages")
			return
		}
	}
}

func (g *GameScene) sendStateToPlayers() {
	for {
		for _, player := range g.state.Players.GetPlayers() {
			if player.GetIpAddress() != "" {
				addr := player.GetIpAddress() + ":" + strconv.Itoa(int(player.GetPort()))
				conn, err := net.Dial("udp", addr)
				if err != nil {
					fmt.Printf("Dial error %v", err)
					return
				}

				gamemsg := utils.CreateGameMessage(1, 1, 2)
				gamemsg.Type = &proto.GameMessage_State{State: &proto.GameMessage_StateMsg{State: g.state}}
				marshal, err := gamemsg.Marshal()
				if err != nil {
					log.Fatal(err)
				}

				_, err = conn.Write(marshal)
				if err != nil {
					log.Fatal("Write failed:", err)
				}

				err = conn.Close()
				if err != nil {
					log.Fatal("Close failed:", err)
				}
			}
		}
		time.Sleep(time.Millisecond * time.Duration(g.state.Config.GetStateDelayMs()))
	}
}

func (g *GameScene) addPlayer(name string, pType proto.PlayerType, view bool, ip string, port int) {
	player := utils.CreatePlayer(name)
	*player.Type = pType
	*player.Port = int32(port)
	*player.IpAddress = ip
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
