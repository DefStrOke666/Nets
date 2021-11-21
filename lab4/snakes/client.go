package snakes

import (
	"fmt"
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"log"
	"net"
)

func (j *JoinScene) receiveAnnouncements() {
	println("Started receiving messages\n")
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	err = l.SetReadBuffer(maxDatagramSize)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, maxDatagramSize)
	for {
		read, addr, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		println("Addr:", addr.String())

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
			msg.Addr = addr
			j.servers = append(j.servers, msg)
			j.serversUpdated = true
		}

		if j.exit {
			println("Stopped receiving announcements\n")
			return
		}
	}
}

func (j *JoinScene) joinServer(addr *net.UDPAddr) bool {
	joinMsg := utils.CreateJoinMessage("client", false)
	conn, err := net.Dial("udp", addr.String())
	if err != nil {
		fmt.Printf("Dial error %v", err)
		return false
	}
	marshal, err := joinMsg.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(marshal)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, maxDatagramSize)
	read, err := conn.Read(b)
	if err != nil {
		log.Fatal("ReadFromUDP failed:", err)
	}

	msg := &proto.GameMessage{}
	err = msg.Unmarshal(b[:read])
	if err != nil {
		log.Fatal(err)
	}

	if msg.Type.Equal(&proto.GameMessage_ErrorMsg{}) {
		errmsg := &proto.GameMessage_ErrorMsg{}
		err = errmsg.Unmarshal(b[:read])
		if err != nil {
			log.Fatal(err)
		}
		println("Connection error:", errmsg.GetErrorMessage())
		return false
	} else if msg.Type.Equal(&proto.GameMessage_AckMsg{}) {
		ackmsg := &proto.GameMessage_AckMsg{}
		err = ackmsg.Unmarshal(b[:read])
		if err != nil {
			log.Fatal(err)
		}
		println("Connected:", ackmsg.String())
		return true
	} else {
		println("Unknown answer")
		return false
	}
}
