package game

import (
	"slimesolver/game/math"
	"strings"
)

type Actor interface {
	Token() Token
	String() string
	GetPosition() math.Vector2
	Transform(g *Game, dir Direction, affectingStates map[Actor]StateChange) (*StateChange, Actor)
	Apply(g *Game, change StateChange)
	Tick(g *Game)
	Solid() bool
	Damage(g *Game)
}

type StateChange struct {
	From    math.Vector2
	Move    math.Vector2
	Updates []Actor // TODO: split this into watching and updating
	Message string
}

func (s StateChange) String() string {
	var sb strings.Builder
	if !s.Move.Equals(math.NegVec) {
		sb.WriteString(s.Move.String())
	}

	if s.Message != "" {
		if sb.Len() > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(s.Message)
	}

	if s.Updates != nil && len(s.Updates) > 0 {
		if sb.Len() > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString("U: [")
		for i, update := range s.Updates {
			sb.WriteString(update.String())
			if i < len(s.Updates)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString("]")
	}

	return sb.String()
}

func (s StateChange) Equals(other StateChange) bool {
	if !s.Move.Equals(other.Move) {
		return false
	}
	if s.Message != other.Message {
		return false
	}
	if len(s.Updates) != len(other.Updates) {
		return false
	}
	for i, update := range s.Updates {
		if update != other.Updates[i] {
			return false
		}
	}
	if !s.From.Equals(other.From) {
		return false
	}

	return s.Move.Equals(other.Move)
}

type PositionComponent struct {
	X, Y int
}

func canMoveTo(g *Game, pos math.Vector2, self Actor) bool {
	if g.GetTokenAt(pos.X, pos.Y) == WallToken {
		return false
	}

	// check
	actors := g.GetActors(pos)
	for _, actor := range actors {
		if actor == self {
			continue
		}

		if actor.Solid() {
			return false
		}
	}

	return true
}

func (p *PositionComponent) GetPosition() math.Vector2 {
	return math.Vector2{p.X, p.Y}
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
