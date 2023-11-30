package game

import "slimesolver/game/math"

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

func (s *Switch) CalculateEdges(g *Game, dir Direction, a Actor) []math.Vector2 {
	if a != nil {
		if a.Token() == BoxToken || a.Token() == SlimeToken {
			// get the doors
			doors := g.GetActorsWithToken(ClosedDoorToken)
			edges := make([]math.Vector2, 0)
			for _, door := range doors {
				edges = append(edges, door.GetPosition())
			}
			return edges
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
