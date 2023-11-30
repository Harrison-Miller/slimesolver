package game

import "slimesolver/game/math"

type Door struct {
	PositionComponent
	open bool
}

func NewDoor(x, y int) *Door {
	return &Door{
		PositionComponent: PositionComponent{x, y},
		open:              false,
	}
}

func (d *Door) Token() Token {
	if d.open {
		return OpenDoorToken
	}
	return ClosedDoorToken
}

func (d *Door) CalculateEdges(g *Game, dir Direction, a Actor) []math.Vector2 {
	if a != nil && a.Token() == SwitchToken {
		return []math.Vector2{{-1, -1}} // dummy edge for opening
	}
	return []math.Vector2{}
}

func (d *Door) ApplyEdges(g *Game, edges []math.Vector2) {
	for _, edge := range edges {
		if edge.X == -1 && edge.Y == -1 {
			d.open = true
			return
		}
	}
	d.open = false
}

func (d *Door) ResolveState(g *Game) {

}

func (d *Door) Solid() bool {
	if d.open {
		return false
	}
	return true
}
