package game

type Spike struct {
	PositionComponent
	up bool
}

func NewSpike(x, y int, up bool) *Spike {
	return &Spike{
		PositionComponent: PositionComponent{x, y},
		up:                up,
	}
}

func (s *Spike) Token() Token {
	if s.up {
		return SpikeUpToken
	}
	return SpikeDownToken
}

func (s *Spike) String() string {
	return string(s.Token())
}

func (s *Spike) Transform(g *Game, dir Direction, affectingStates map[Actor]StateChange) (*StateChange, Actor) {
	return nil, nil
}

func (s *Spike) Apply(g *Game, change StateChange) {

}

func (s *Spike) Tick(g *Game) {
	s.up = !s.up

	if s.up {
		actors := g.GetActors(s.GetPosition())
		for _, actor := range actors {
			if actor != s {
				actor.Damage(g)
			}
		}
	}
}

func (s *Spike) Solid() bool {
	return false
}

func (s *Spike) Damage(g *Game) {

}
