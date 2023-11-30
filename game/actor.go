package game

import (
	"fmt"
	"slimesolver/game/math"
)

type Actor interface {
	Token() Token
	GetPosition() math.Vector2
	CalculateEdges(g *Game, dir Direction, a Actor) []math.Vector2
	ApplyEdges(g *Game, edges []math.Vector2)
	ResolveState(g *Game)
	Solid() bool
}

type PositionComponent struct {
	X, Y int
}

func (p *PositionComponent) GetPosition() math.Vector2 {
	return math.Vector2{p.X, p.Y}
}

type MoveableComponent struct {
	PositionComponent
}

func (m *MoveableComponent) ApplyEdges(g *Game, edges []math.Vector2) {
	for _, edge := range edges {
		if g.IsWallOrEdge(edge.X, edge.Y) {
			fmt.Println("Can't move because of wall or edge")
			continue
		}

		actors := g.GetActors(edge.X, edge.Y)
		foundSolid := false
		for _, actor := range actors {
			if actor.Solid() {
				fmt.Println("Can't move because of solid actor")
				foundSolid = true
				break
			}
		}

		if foundSolid {
			continue
		}

		m.X = edge.X
		m.Y = edge.Y
	}
}

func moveVector(vec math.Vector2, dir Direction) math.Vector2 {
	switch dir {
	case Up:
		return math.Vector2{vec.X, vec.Y - 1}
	case Down:
		return math.Vector2{vec.X, vec.Y + 1}
	case Left:
		return math.Vector2{vec.X - 1, vec.Y}
	case Right:
		return math.Vector2{vec.X + 1, vec.Y}
	default:
		return vec
	}
}

func directionBetween(a, b math.Vector2) Direction {
	if a.X == b.X {
		if a.Y > b.Y {
			return Up
		} else {
			return Down
		}
	} else {
		if a.X > b.X {
			return Left
		} else {
			return Right
		}
	}
}
