package game

import (
	"fmt"
	"slimesolver/game/math"
)

type Box struct {
	PositionComponent
}

func NewBox(x, y int) *Box {
	return &Box{
		PositionComponent: PositionComponent{x, y},
	}
}

func (b *Box) Token() Token {
	return BoxToken
}

func (b *Box) CalculateEdges(g *Game, dir Direction, a Actor) []math.Vector2 {
	if a != nil {
		if a.Token() == SlimeToken {
			move := moveVector(b.GetPosition(), dir)
			// check if we can move
			if g.IsWallOrEdge(move.X, move.Y) {
				return []math.Vector2{b.GetPosition()}
			}

			return []math.Vector2{move}
		}
		return []math.Vector2{}
	}
	return []math.Vector2{b.GetPosition()}
}

func (b *Box) ApplyEdges(g *Game, edges []math.Vector2) {
	for _, edge := range edges {
		if g.IsWallOrEdge(edge.X, edge.Y) {
			fmt.Println("Can't move because of wall or edge")
			continue
		}

		actors := g.GetActors(edge.X, edge.Y)
		foundSolid := false
		for _, actor := range actors {
			if actor == b {
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

		b.X = edge.X
		b.Y = edge.Y
	}
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
