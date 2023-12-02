package game

import (
	"slimesolver/game/math"
)

type Door struct {
	PositionComponent
	open bool
}

func NewDoor(x, y int) *Door {
	return &Door{
		PositionComponent: PositionComponent{x, y},
	}
}

func (d *Door) Token() Token {
	if d.open {
		return OpenDoorToken
	}
	return ClosedDoorToken
}

func (d *Door) String() string {
	return string(d.Token())
}

func (d *Door) Transform(g *Game, dir Direction, affectingStates map[Actor]StateChange) (*StateChange, Actor) {
	open := false
	var openActor Actor
	var otherActor Actor
	for actor, state := range affectingStates {
		if actor.Token() == SwitchToken {
			open = true
			openActor = actor
		}

		if actor.Token() == BoxToken || actor.Token() == SlimeToken {
			if state.Move.Equals(d.GetPosition()) {
				otherActor = actor
			}
		}
	}

	if open {
		if otherActor != nil { // doing this forces the slime to move after
			openActor = otherActor
		}
		return &StateChange{
			Move:    math.NegVec,
			Message: "open",
		}, openActor
	}

	return &StateChange{
		Move:    math.NegVec,
		Message: "close",
	}, otherActor
}

func (d *Door) Apply(g *Game, change StateChange) {
	d.open = change.Message == "open"
}

func (d *Door) Tick(g *Game) {
	if !d.open {
		actors := g.GetActors(d.GetPosition())
		for _, actor := range actors {
			if actor != d {
				g.Kill(actor)
			}
		}
	}
}

func (d *Door) Solid() bool {
	return !d.open
}

func (d *Door) Damage(g *Game) {

}
