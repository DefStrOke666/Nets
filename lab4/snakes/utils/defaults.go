package utils

import (
	"fmt"
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/spf13/viper"
	"log"
	"math/rand"
)

type Players struct {
	AdminName  string
	PlayerName string
}

type Game struct {
	Width         int32
	Height        int32
	FoodStatic    int32
	FoodPerPlayer float32
	StateDelayMs  int32
	DeadFoodProb  float32
	PingDelayMs   int32
	NodeTimeoutMs int32
}

type Config struct {
	Game        Game
	PlayerNames Players
}

var Conf Config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatalf("failed to load: %v", err)
	}

	fmt.Printf("Config: %+v\n", Conf)
}

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
	*conf.Width = Conf.Game.Width
	*conf.Height = Conf.Game.Height
	*conf.FoodStatic = Conf.Game.FoodStatic
	*conf.FoodPerPlayer = Conf.Game.FoodPerPlayer
	*conf.StateDelayMs = Conf.Game.StateDelayMs
	*conf.DeadFoodProb = Conf.Game.DeadFoodProb
	*conf.PingDelayMs = Conf.Game.PingDelayMs
	*conf.NodeTimeoutMs = Conf.Game.NodeTimeoutMs
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

func CreateJoin(name string, view bool) *proto.GameMessage_Join {
	joinMsg := CreateJoinMessage(name, view)
	join := &proto.GameMessage_Join{Join: joinMsg}

	return join
}

func CreateJoinMessage(name string, view bool) *proto.GameMessage_JoinMsg {
	joinMsg := &proto.GameMessage_JoinMsg{
		PlayerType: new(proto.PlayerType),
		OnlyView:   new(bool),
		Name:       new(string),
	}

	*joinMsg.PlayerType = proto.PlayerType_HUMAN
	*joinMsg.OnlyView = view
	*joinMsg.Name = name

	return joinMsg
}

func CreateGameMessage(seq int64, senderId, receiverId int32) *proto.GameMessage {
	gameMsh := &proto.GameMessage{
		MsgSeq:     new(int64),
		SenderId:   new(int32),
		ReceiverId: new(int32),
	}

	*gameMsh.MsgSeq = seq
	*gameMsh.SenderId = senderId
	*gameMsh.ReceiverId = receiverId

	return gameMsh
}
