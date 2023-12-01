package game

import (
	"fmt"
	"slimesolver/game/math"
)

type Slime struct {
	PositionComponent
}

func NewSlime(x, y int) *Slime {
	return &Slime{
		PositionComponent: PositionComponent{x, y},
	}
}

func (s *Slime) Token() Token {
	return SlimeToken
}

func (s *Slime) CalculateEdges(g *Game, dir Direction, a Actor) []math.Vector2 {
	if a == nil && dir != Zero {
		move := moveVector(s.GetPosition(), dir)

		// check if we can move
		if g.IsWallOrEdge(move.X, move.Y) {
			return []math.Vector2{s.GetPosition()}
		}

		return []math.Vector2{move}
	}
	return []math.Vector2{}
}

func (s *Slime) ApplyEdges(g *Game, edges []math.Vector2) {
	for _, edge := range edges {
		if g.IsWallOrEdge(edge.X, edge.Y) {
			fmt.Println("Can't move because of wall or edge")
			continue
		}

		actors := g.GetActors(edge.X, edge.Y)
		foundSolid := false
		for _, actor := range actors {
			if actor == s {
				continue
			}

			if actor.Solid() {
				fmt.Printf("Can't move because of solid actor - %s\n", string(actor.Token()))
				foundSolid = true
				break
			}
		}

		if foundSolid {
			continue
		}

		s.X = edge.X
		s.Y = edge.Y
	}
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
