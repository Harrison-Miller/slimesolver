package game

type Pusher struct {
	PositionComponent
	active bool
}

func NewPusher(x, y int, active bool) *Pusher {
	return &Pusher{
		PositionComponent: PositionComponent{x, y},
		active:            active,
	}
}

func (p Pusher) Token() Token {
	if p.active {
		return PusherActiveToken
	}
	return PusherToken
}

func (p Pusher) String() string {
	return string(p.Token())
}

func (p Pusher) Transform(g *Game, dir Direction, affectingStates map[Actor]StateChange) (*StateChange, Actor) {
	return nil, nil
}

func (p Pusher) Apply(g *Game, change StateChange) {

}

func (p Pusher) Tick(g *Game) {
	p.active = !p.active
}

func (p Pusher) Solid() bool {
	return true
}

func (p Pusher) Damage(g *Game) {

}
