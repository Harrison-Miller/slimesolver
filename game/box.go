package game

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

func (b *Box) String() string {
	return string(b.Token())
}

func (b *Box) Transform(g *Game, dir Direction, affectingStates map[Actor]StateChange) (*StateChange, Actor) {
	for actor, change := range affectingStates {
		if actor.Token() == SlimeToken {
			if change.Move.Equals(b.GetPosition()) {
				dir = directionBetween(change.From, b.GetPosition())
				move := moveVector(b.GetPosition(), dir)
				return &StateChange{
					Move: move,
				}, actor
			}
		}
	}

	return &StateChange{
		Move: b.GetPosition(),
	}, nil
}

func (b *Box) Apply(g *Game, change StateChange) {
	// check if we can move
	canMove := canMoveTo(g, change.Move, b)
	if !canMove {
		return
	}

	b.X = change.Move.X
	b.Y = change.Move.Y
}

func (b *Box) Tick(g *Game) {
	// fall into pit
	if g.IsPit(b.X, b.Y) {
		g.SetTokenAt(b.X, b.Y, EmptyToken)
		g.Kill(b)
		return
	}
}

func (b *Box) Solid() bool {
	return true
}

func (b *Box) Damage(g *Game) {

}
