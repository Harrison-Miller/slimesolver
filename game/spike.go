package game

import "slimesolver/game/math"

type Spike struct {
	PositionComponent
	up bool
}

func NewSpike(x, y int, up bool) *Spike {
	return &Spike{
		PositionComponent: PositionComponent{x, y},
		up:                up,
	}
}

func (s *Spike) Token() Token {
	if s.up {
		return SpikeUpToken
	}
	return SpikeDownToken
}

func (s *Spike) CalculateEdges(g *Game, dir Direction, a Actor) []math.Vector2 {
	return []math.Vector2{}
}

func (s *Spike) ApplyEdges(g *Game, edges []math.Vector2) {

}

func (s *Spike) ResolveState(g *Game) {
	s.up = !s.up
	if s.up {
		actors := g.GetActors(s.GetPosition())
		for _, actor := range actors {
			if actor.Token() == SlimeToken {
				g.RemoveActor(actor)
			}
		}
	}
}

func (s *Spike) Solid() bool {
	return false
}
