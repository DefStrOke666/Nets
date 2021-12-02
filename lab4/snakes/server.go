package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"log"
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
			return
		}

		if g.exit {
			println("Stopped listening for messages")
			return
		}
	}
}
