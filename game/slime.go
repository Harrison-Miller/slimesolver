package game

import (
	"fmt"
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

func possibleBlockerStates(affectingStates AffectingStates) map[Actor]StateChange {
	possibleBlockers := affectingStates.GoingToStates
	for actor, change := range affectingStates.WatchingStates {
		possibleBlockers[actor] = change
	}
	return possibleBlockers
}

func (s *Slime) Transform(g *Game, dir Direction, affectingStates AffectingStates) (*StateChange, Actor) {
	pos := s.GetPosition()
	move := moveVector(pos, dir)
	nextChange := &StateChange{
		Move: move,
	}
	var parent Actor

	// things moving to where we are need to wait for us to move
	// so they become our parent
	for actor, _ := range affectingStates.OnToStates {
		parent = actor
	}

	// can't move if we're going to hit a wall
	if g.IsWallOrEdge(move.X, move.Y) {
		nextChange.Move = pos
	}

	// doors that aren't opening block our movement
	if len(affectingStates.GoingToStates) > 0 {
		fmt.Println("going to states: ", affectingStates.GoingToStates)
	}
	possibleBlockers := possibleBlockerStates(affectingStates) // includes going to and watched states
	for actor, change := range possibleBlockers {
		token := actor.Token()
		if token == OpenDoorToken || token == ClosedDoorToken {
			if change.Message == "close" {
				nextChange.Move = pos
				nextChange.Watching = append(nextChange.Watching, actor) // we need to keep track of this actor in future transforms
			}
		}
	}

	if s.small {
		// small slime is blocked by boxes and normal slimes
		for actor, _ := range possibleBlockers {
			token := actor.Token()
			switch token {
			case BoxToken, SlimeToken:
				nextChange.Move = pos
				nextChange.Watching = append(nextChange.Watching, actor) // we need to keep track of this actor in future transforms
			}
		}

		// if this slime is not going to move
		// and another small slime will move into it this turn
		// then we need to grow
		if nextChange.Move.Equals(pos) {
			for actor, _ := range affectingStates.OnToStates {
				if actor.Token() == SmallSlimeToken {
					nextChange.Message = "grow"
				}
			}
		}

		// if we're moving into another small slime, we need to combine
		for actor, change := range affectingStates.GoingToStates {
			if actor.Token() == SmallSlimeToken && change.Message == "grow" {
				nextChange.Message = "combine"
			}
		}
	}

	return nextChange, parent
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
