package game

import (
	"github.com/borodun/nsu-nets/lab4/game/proto"
	"github.com/hajimehoshi/ebiten/v2"
)

type ServerList struct {
	servers []proto.GameMessage_AnnouncementMsg

	serverImg         []*Picture
	serversBackground *ebiten.Image
	infoImg           *ebiten.Image
	infoBackground    *ebiten.Image

	listW, listH int
	infoW, infoH int

	selectedServer int
}

func generateConfig() *proto.GameConfig {
	conf := &proto.GameConfig{}
	conf.Width = new(int32)
	conf.Height = new(int32)
	conf.FoodStatic = new(int32)
	conf.FoodPerPlayer = new(float32)
	conf.StateDelayMs = new(int32)
	conf.DeadFoodProb = new(float32)
	conf.PingDelayMs = new(int32)
	conf.NodeTimeoutMs = new(int32)

	*conf.Width = 20
	*conf.Height = 20
	*conf.FoodStatic = 5
	*conf.FoodPerPlayer = 1
	*conf.StateDelayMs = 5
	*conf.DeadFoodProb = 20
	*conf.PingDelayMs = 20
	*conf.NodeTimeoutMs = 200
	return conf
}

func generatePlayer() *proto.GamePlayer {
	player := &proto.GamePlayer{}
	player.Name = new(string)
	*player.Name = "Player"
	return player
}

func generatePlayers(count int) *proto.GamePlayers {
	players := &proto.GamePlayers{}
	for i := 0; i < count; i++ {
		players.Players = append(players.Players, generatePlayer())
	}
	return players
}

func generateServers(count int) []proto.GameMessage_AnnouncementMsg {
	servers := make([]proto.GameMessage_AnnouncementMsg, count)
	for i := range servers {
		servers[i].Players = generatePlayers(5)
		servers[i].CanJoin = new(bool)
		*servers[i].CanJoin = true
		servers[i].Config = generateConfig()
	}

	return servers
}
