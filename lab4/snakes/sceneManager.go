package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	transitionFrom = ebiten.NewImage(screenWidth, screenHeight)
	transitionTo   = ebiten.NewImage(screenWidth, screenHeight)
)

type Scene interface {
	Update(state *GameState) error
	Draw(screen *ebiten.Image)
}

const transitionMaxCount = 15

type SceneManager struct {
	current         Scene
	next            Scene
	transitionCount int
}

type GameState struct {
	SceneManager *SceneManager
	State        *proto.GameState
}

func (s *SceneManager) Update() error {
	if s.transitionCount == 0 {
		return s.current.Update(&GameState{
			SceneManager: s,
		})
	}

	s.transitionCount--
	if s.transitionCount > 0 {
		return nil
	}

	s.current = s.next
	s.next = nil
	return nil
}

func (s *SceneManager) Draw(r *ebiten.Image) {
	if s.transitionCount == 0 {
		s.current.Draw(r)
		return
	}

	transitionFrom.Clear()
	s.current.Draw(transitionFrom)

	transitionTo.Clear()
	s.next.Draw(transitionTo)

	r.DrawImage(transitionFrom, nil)

	alpha := 1 - float64(s.transitionCount)/float64(transitionMaxCount)
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, alpha)
	r.DrawImage(transitionTo, op)
}

func (s *SceneManager) GoTo(scene Scene) {
	if s.current == nil {
		s.current = scene
	} else {
		s.next = scene
		s.transitionCount = transitionMaxCount
	}
}
