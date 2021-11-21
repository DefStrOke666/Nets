package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
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
