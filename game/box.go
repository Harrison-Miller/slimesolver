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

func (b *Box) Transform(g *Game, dir Direction, affectingStates AffectingStates) (*StateChange, Actor) {
	pos := b.GetPosition()
	// something pushing the box
	for actor, change := range affectingStates.OnToStates {
		token := actor.Token()
		switch token {
		case SlimeToken:
			dir = directionBetween(change.From, b.GetPosition())
			move := moveVector(pos, dir)
			return &StateChange{
				Move: move,
			}, actor
		}
	}

	// otherwise slime stands still
	return &StateChange{
		Move: pos,
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
