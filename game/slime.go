package game

import (
	"slimesolver/game/math"
)

type Slime struct {
	MoveableComponent
}

func NewSlime(x, y int) *Slime {
	return &Slime{
		MoveableComponent: MoveableComponent{
			PositionComponent: PositionComponent{x, y},
		},
	}
}

func (s *Slime) Token() Token {
	return SlimeToken
}

func (s *Slime) CalculateEdges(g *Game, dir Direction, a Actor) []math.Vector2 {
	if a == nil && dir != Zero {
		return []math.Vector2{moveVector(s.GetPosition(), dir)}
	}
	return []math.Vector2{}
}

func (s *Slime) ResolveState(g *Game) {
	token := g.GetTokenAt(s.X, s.Y)
	if token == PitToken {
		g.RemoveActor(s)
	}
}

func (s *Slime) Solid() bool {
	return true
}
