package snakes

import (
	"fmt"
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"net"
	"strings"
	"time"
)

func (j *JoinScene) receiveAnnouncements() {
	println("Started receiving messages\n")
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatal("ResolveUDPAddr:", err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("ListenMulticastUDP:", err)
	}
	err = l.SetReadBuffer(maxDatagramSize)
	if err != nil {
		log.Fatal("SetReadBuffer:", err)
	}

	b := make([]byte, maxDatagramSize)
	for {
		read, addr, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		println("Announcement from addr:", addr.String())

		msg := &proto.GameMessage_AnnouncementMsg{}
		err = msg.Unmarshal(b[:read])
		if err != nil {
			log.Fatal(err)
		}

		serverExists := false
		for _, server := range j.servers {
			if msg.Equal(server) {
				serverExists = true
			}
		}
		if !serverExists {
			j.servers = append(j.servers, msg)
			j.serversUpdated = true
		}

		if j.exit {
			println("Stopped receiving announcements\n")
			return
		}
	}
}

func (g *GameScene) joinServer(view bool) bool {
	conn, err := net.Dial("udp", g.servAddr)
	if err != nil {
		fmt.Printf("Dial error %v", err)
		return false
	}
	println("Connected to:", conn.RemoteAddr().String(), " server addr:", g.servAddr)

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Conn close:", err)
		}
	}(conn)

	g.port = strings.Split(conn.LocalAddr().String(), ":")[1]
	println("My port:", g.port)

	gamemsg := utils.CreateGameMessage(1, 2, 1)
	gamemsg.Type = utils.CreateJoin(g.playerName, view)
	marshal, err := gamemsg.Marshal()
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write(marshal)
	if err != nil {
		log.Fatal("Write failed:", err)
	}

	b := make([]byte, maxDatagramSize)
	read, err := conn.Read(b)
	if err != nil {
		log.Fatal("Read failed:", err)
	}

	msg := &proto.GameMessage{}
	err = msg.Unmarshal(b[:read])
	if err != nil {
		log.Fatal(err)
	}

	switch message := msg.Type.(type) {
	case *proto.GameMessage_Ack:
		println("Ack:", message.Ack.String())
		return true
	case *proto.GameMessage_Error:
		println("Error:", message.Error.GetErrorMessage())
		return false
	default:
		println("Unknown answer")
		return false
	}
}

func (g *GameScene) receiveMessages() {
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
		read, _, err := conn.ReadFromUDP(b)
		if err != nil {
			log.Print("ReadFromUDP failed:", err)
			continue
		}
		//println("Connection from:", addr.String())

		message := &proto.GameMessage{}
		err = message.Unmarshal(b[:read])
		if err != nil {
			log.Print("Unmarshall error:", err)
			continue
		}

		switch msg := message.Type.(type) {
		case *proto.GameMessage_State:
			//println("State from addr:", addr.String())
			state := msg.State
			if g.state.GetStateOrder() < state.State.GetStateOrder() {
				g.state = state.State
				//println("State:", state.String())
				g.drawScore()
			} else {
				println("No new state")
			}
			continue
		}

		if g.exit {
			println("Stopped listening for messages")
			return
		}
		time.Sleep(time.Millisecond * time.Duration(g.state.Config.GetStateDelayMs()))
	}
}

func (g *GameScene) sendDirection() {
	direction := proto.Direction(0)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		direction = proto.Direction_UP
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		direction = proto.Direction_DOWN
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		direction = proto.Direction_RIGHT
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		direction = proto.Direction_LEFT
	}

	conn, err := net.Dial("udp", g.servAddr)
	if err != nil {
		fmt.Printf("Dial error %v", err)
		return
	}
	println("Connected to:", conn.RemoteAddr().String())

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Conn close:", err)
		}
	}(conn)

	gamemsg := utils.CreateGameMessage(1, 2, 1)
	gamemsg.Type = &proto.GameMessage_Steer{Steer: &proto.GameMessage_SteerMsg{Direction: &direction}}
	marshal, err := gamemsg.Marshal()
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write(marshal)
	if err != nil {
		log.Fatal("Write failed:", err)
	}
}
