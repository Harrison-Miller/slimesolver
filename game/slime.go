package game

import (
	"slimesolver/game/math"
)

type Slime struct {
	PositionComponent
	small        bool
	lastPosition math.Vector2
}

func NewSlime(x, y int, small bool) *Slime {
	return &Slime{
		PositionComponent: PositionComponent{x, y},
		small:             small,
	}
}

func (s *Slime) Token() Token {
	if s.small {
		return SmallSlimeToken
	}
	return SlimeToken
}

func (s *Slime) String() string {
	return string(s.Token())
}

func (s *Slime) Transform(g *Game, dir Direction, affectingStates map[Actor]StateChange) (*StateChange, Actor) {
	pos := s.GetPosition()
	move := moveVector(pos, dir)
	canMove := !g.IsWallOrEdge(move.X, move.Y)
	updates := make([]Actor, 0)
	var parent Actor
	message := ""
	for a, change := range affectingStates {
		if change.Move.Equals(s.GetPosition()) {
			parent = a
		}

		if change.Message == "grow" && s.small {
			message = "combine"
		}

		// what we're moving into
		apos := a.GetPosition()
		if apos.Equals(move) {
			token := a.Token()
			if token == OpenDoorToken || token == ClosedDoorToken {
				if change.Message == "close" {
					canMove = false
					updates = append(updates, a)
				}
			}

			if (token == BoxToken || token == SlimeToken) && s.small {
				if change.Move.Equals(apos) {
					canMove = false
					updates = append(updates, a)
				}
			}
		}
	}

	if !canMove {
		if parent != nil && s.small && parent.Token() == SmallSlimeToken {
			message = "grow"
			updates = append(updates, parent)
		}

		return &StateChange{
			Move:    pos,
			Message: message,
			Updates: updates,
		}, parent
	}

	return &StateChange{
		Move:    move,
		Message: message,
	}, parent
}

func (s *Slime) Apply(g *Game, change StateChange) {
	if change.Message == "grow" {
		s.small = false
	} else if change.Message == "combine" {
		g.Kill(s)
	}

	// check if we can move
	canMove := canMoveTo(g, change.Move, s)
	if !canMove {
		return
	}

	s.lastPosition = s.GetPosition()
	s.X = change.Move.X
	s.Y = change.Move.Y
}

func (s *Slime) Tick(g *Game) {
	// die if we're on a pit
	if g.IsPit(s.X, s.Y) {
		g.Kill(s)
		return
	}
}

func (s *Slime) Solid() bool {
	return true
}

func (s *Slime) Damage(g *Game) {
	if s.small {
		g.Kill(s)
		return
	}

	s.small = true

	spawnLocations := s.getSpawnLocations()
	for _, loc := range spawnLocations {
		if canMoveTo(g, loc, s) {
			g.AddActor(NewSlime(loc.X, loc.Y, true))
			return
		}
	}
}

func (s *Slime) getSpawnLocations() []math.Vector2 {
	pos := s.GetPosition()
	return []math.Vector2{
		s.lastPosition,
		moveVector(pos, Up),
		moveVector(pos, Down),
		moveVector(pos, Left),
		moveVector(pos, Right),
	}
}
