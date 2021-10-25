package main

import (
	"fmt"
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	srvAddr         = "239.192.0.4:9192"
	maxDatagramSize = 8192
)

func main() {
	ping(srvAddr)
}

func ping(a string) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	for {
		msg := proto.NewPopulatedGameMessage_AnnouncementMsg(rand.New(rand.NewSource(time.Now().UnixNano())), false)
		marshal, err := msg.Marshal()
		if err != nil {
			log.Fatal(err)
		}
		c.Write(marshal)
		fmt.Println("Wrote: " + msg.String())
		time.Sleep(5 * time.Second)
	}
}
