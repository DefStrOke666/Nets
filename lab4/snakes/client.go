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

func (j *JoinScene) joinServer(addr *net.UDPAddr, view bool) bool {
	joinMsg := utils.CreateJoinMessage("client", view)
	conn, err := net.Dial("udp", addr.String())
	if err != nil {
		fmt.Printf("Dial error %v", err)
		return false
	}
	println("Connected to:", conn.RemoteAddr().String())

	marshal, err := joinMsg.Marshal()
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
