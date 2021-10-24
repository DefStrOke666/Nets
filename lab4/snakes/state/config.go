package state

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"math/rand"
)

func NewDefaultGameConfig() *proto.GameConfig {
	conf := &proto.GameConfig{
		Width:         new(int32),
		Height:        new(int32),
		FoodStatic:    new(int32),
		FoodPerPlayer: new(float32),
		StateDelayMs:  new(int32),
		DeadFoodProb:  new(float32),
		PingDelayMs:   new(int32),
		NodeTimeoutMs: new(int32),
	}
	*conf.Width = 30
	*conf.Height = 15
	*conf.FoodStatic = 2
	*conf.FoodPerPlayer = 1
	*conf.StateDelayMs = 100
	*conf.DeadFoodProb = 0.1
	*conf.PingDelayMs = 50
	*conf.NodeTimeoutMs = 1000
	return conf
}

func CreatePlayer(name string) *proto.GamePlayer {
	player := &proto.GamePlayer{
		Name:      new(string),
		Id:        new(int32),
		IpAddress: new(string),
		Port:      new(int32),
		Role:      new(proto.NodeRole),
		Type:      new(proto.PlayerType),
		Score:     new(int32),
	}

	*player.Name = name
	*player.Id = rand.Int31n(100)
	*player.IpAddress = ""
	*player.Port = 10000
	*player.Role = proto.NodeRole_NORMAL
	*player.Type = proto.PlayerType_HUMAN
	*player.Score = 0

	return player
}

func CreateSnake(id int, head *proto.GameState_Coord) *proto.GameState_Snake {
	snake := &proto.GameState_Snake{
		PlayerId:      new(int32),
		Points:        make([]*proto.GameState_Coord, 1),
		State:         new(proto.GameState_Snake_SnakeState),
		HeadDirection: new(proto.Direction),
	}

	*snake.PlayerId = int32(id)
	snake.Points[0] = head
	*snake.State = proto.GameState_Snake_ALIVE
	*snake.HeadDirection = proto.Direction(rand.Intn(4) + 1)

	return snake
}

func CreateCoord(x, y int) *proto.GameState_Coord {
	coord := &proto.GameState_Coord{
		X: new(int32),
		Y: new(int32),
	}

	*coord.X = int32(x)
	*coord.Y = int32(y)

	return coord
}
