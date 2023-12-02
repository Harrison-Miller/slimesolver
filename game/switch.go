package game

import (
	"slimesolver/game/math"
)

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

func (s *Switch) String() string {
	return string(s.Token())
}

func getDoors(g *Game) []Actor {
	return g.GetActorsWithTokens([]Token{ClosedDoorToken, OpenDoorToken})
}

func (s *Switch) Transform(g *Game, dir Direction, affectingStates map[Actor]StateChange) (*StateChange, Actor) {
	if len(affectingStates) == 0 {
		return nil, nil
	}

	doors := getDoors(g)

	pressed := false
	var pressingActor Actor
	for actor, change := range affectingStates {
		if change.Move.Equals(s.GetPosition()) {
			if actor.Token() == SlimeToken || actor.Token() == BoxToken ||
				(actor.Token() == SmallSlimeToken && change.Message == "grow") {
				pressed = true
				pressingActor = actor
				break
			}
		}
	}

	if pressed {
		return &StateChange{
			Move:    math.NegVec,
			Updates: doors,
		}, pressingActor
	}

	// check if something not moving is on top of us
	pos := s.GetPosition()
	actors := g.GetActors(pos)
	for _, actor := range actors {
		if actor.Token() == SlimeToken || actor.Token() == BoxToken {
			// check if the slime/box is moving away from the switch
			if state, ok := affectingStates[actor]; ok {
				if state.From.Equals(pos) {
					continue
				}
			}

			return &StateChange{
				Move:    math.NegVec,
				Updates: doors,
			}, nil // we can't pass a parent actor because it's not moving (doesn't have its own state change to connect to).
		}
	}

	return nil, nil
}

func (s *Switch) Apply(g *Game, change StateChange) {
}

func (s *Switch) Tick(g *Game) {
}

func (s *Switch) Solid() bool {
	return false
}

func (s *Switch) Damage(g *Game) {

}
