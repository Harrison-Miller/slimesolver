package game

import (
	"slimesolver/game/math"
)

type Box struct {
	MoveableComponent
}

func NewBox(x, y int) *Box {
	return &Box{
		MoveableComponent: MoveableComponent{
			PositionComponent: PositionComponent{x, y},
		},
	}
}

func (b *Box) Token() Token {
	return BoxToken
}

func (b *Box) CalculateEdges(g *Game, dir Direction, a Actor) []math.Vector2 {
	if a != nil && a.Token() == SlimeToken {
		return []math.Vector2{moveVector(b.GetPosition(), dir)}
	}
	return []math.Vector2{}
}

func (b *Box) ResolveState(g *Game) {
	token := g.GetTokenAt(b.X, b.Y)
	if token == PitToken {
		g.RemoveActor(b)
		g.SetTokenAt(b.X, b.Y, EmptyToken)
	}
}

func (b *Box) Solid() bool {
	return true
}
