package game

import "slimesolver/game/math"

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

func (d *Door) Transform(g *Game, dir Direction, affectingStates AffectingStates) (*StateChange, Actor) {
	nextChange := &StateChange{
		Move:    math.NegVec,
		Message: "close",
	}
	var parent Actor

	// switch activated by something
	for actor, _ := range affectingStates.UpdateStates {
		token := actor.Token()
		switch token {
		case SwitchToken:
			nextChange.Message = "open"
			parent = actor
		}
	}

	// make moving objects dependent on the door
	// parents move after children
	for actor, _ := range affectingStates.OnToStates {
		parent = actor
	}

	return nextChange, parent
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
