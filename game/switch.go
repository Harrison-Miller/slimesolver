package game

import (
	"slimesolver/game/math"
)

type Switch struct {
	PositionComponent
}

func NewSwitch(x, y int) *Switch {
	return &Switch{
		PositionComponent: PositionComponent{x, y},
	}
}

func (s *Switch) Token() Token {
	return SwitchToken
}

func doorEdges(g *Game) []math.Vector2 {
	doors := g.GetActorsWithTokens([]Token{ClosedDoorToken, OpenDoorToken})
	edges := make([]math.Vector2, 0)
	for _, door := range doors {
		edges = append(edges, door.GetPosition())
	}

	return edges
}

func (s *Switch) CalculateEdges(g *Game, dir Direction, a Actor) []math.Vector2 {
	if a != nil {
		// if something is moving onto us
		if a.Token() == BoxToken || a.Token() == SlimeToken {
			return doorEdges(g)
		}
	}

	return []math.Vector2{}
}

func (s *Switch) ApplyEdges(g *Game, edges []math.Vector2) {

}

func (s *Switch) ResolveState(g *Game) {

}

func (s *Switch) Solid() bool {
	return false
}
