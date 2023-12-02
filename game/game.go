package game

import (
	"fmt"
	"slimesolver/game/math"
	"sort"
	"strings"
)

type Token rune

const (
	WallToken       Token = '#'
	EmptyToken      Token = '.'
	PitToken        Token = 'O'
	SlimeToken      Token = '@'
	SmallSlimeToken       = 'o'
	BoxToken        Token = 'B'
	SwitchToken     Token = 'x'
	ClosedDoorToken Token = 'D'
	OpenDoorToken   Token = '_'
	SpikeUpToken    Token = '^'
	SpikeDownToken  Token = '-'
)

type Direction int

const (
	Zero Direction = iota
	Up
	Down
	Left
	Right
)

func dirString(dir Direction) string {
	switch dir {
	case Up:
		return "up"
	case Down:
		return "down"
	case Left:
		return "left"
	case Right:
		return "right"
	default:
		return "zero"
	}
}

type Game struct {
	// 2d array representing the static part of the board
	// y, x
	board [][]Token

	// there can be multiple actors on the same tile
	actors []Actor

	killQueue []Actor

	logging bool
}

func NewGame(logging bool) *Game {
	return &Game{
		logging: logging,
	}
}

func (g *Game) Println(args ...interface{}) {
	if g.logging {
		fmt.Println(args...)
	}
}

func (g *Game) Printf(format string, args ...interface{}) {
	if g.logging {
		fmt.Printf(format, args...)
	}
}

func (g *Game) GetTokenAt(x, y int) Token {
	return g.board[y][x]
}

func (g *Game) SetTokenAt(x, y int, token Token) {
	g.board[y][x] = token
}

func (g *Game) IsWallOrEdge(x, y int) bool {
	if x < 0 || y < 0 || y >= len(g.board) || x >= len(g.board[y]) {
		return true
	}

	return g.GetTokenAt(x, y) == WallToken
}

func (g *Game) IsPit(x, y int) bool {
	return g.GetTokenAt(x, y) == PitToken
}

func (g *Game) AddActor(actor Actor) {
	g.actors = append(g.actors, actor)
}

func (g *Game) GetActors(pos math.Vector2) []Actor {
	l := make([]Actor, 0)
	for _, entity := range g.actors {
		v := entity.GetPosition()
		if pos.Equals(v) {
			l = append(l, entity)
		}
	}
	return l
}

func (g *Game) GetActorsWithTokens(tokens []Token) []Actor {
	l := make([]Actor, 0)
	for _, actor := range g.actors {
		for _, token := range tokens {
			if actor.Token() == token {
				l = append(l, actor)
				break
			}
		}
	}
	return l
}

func (g *Game) Kill(actor Actor) {
	g.killQueue = append(g.killQueue, actor)
}

func (g *Game) RemoveActor(actor Actor) {
	for i, e := range g.actors {
		if e == actor {
			g.actors = append(g.actors[:i], g.actors[i+1:]...)
			return
		}
	}
}

func getPriorityToken(actors []Actor) Token {
	var token Token
	priority := -1
	for _, actor := range actors {
		t := actor.Token()
		p := 0
		switch t {
		case SlimeToken:
			p = 10
		case BoxToken:
			p = 5
		default:
			p = 0
		}
		if p > priority {
			token = t
			priority = p
		}
	}
	return token
}

func (g *Game) String() string {
	var sb strings.Builder
	for y, row := range g.board {
		for x, token := range row {
			actors := g.GetActors(math.Vector2{x, y})
			if len(actors) > 0 {
				token = getPriorityToken(actors)
			}
			sb.WriteRune(rune(token))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func cleanState(state string) string {
	state = strings.ReplaceAll(state, " ", "")
	state = strings.ReplaceAll(state, "\t", "")
	// replace line ending
	state = strings.ReplaceAll(state, "\r\n", "\n")
	return state
}

func (g *Game) Parse(state string) error {
	state = cleanState(state)

	lines := strings.Split(state, "\n")
	width := len(lines[0])
	height := len(lines)

	// initialize board
	g.board = make([][]Token, height)
	for i := range g.board {
		g.board[i] = make([]Token, width)
	}

	// initialize actors
	g.actors = make([]Actor, 0)

	for y, line := range lines {
		if len(line) != width {
			return fmt.Errorf("invalid line length: %d", len(line))
		}

		// construct board and objects
		for x, c := range line {
			switch Token(c) {
			case WallToken:
				g.board[y][x] = WallToken
			case EmptyToken:
				g.board[y][x] = EmptyToken
			case PitToken:
				g.board[y][x] = PitToken
			case SlimeToken, SmallSlimeToken:
				g.board[y][x] = EmptyToken
				g.actors = append(g.actors, NewSlime(x, y, Token(c) == SmallSlimeToken))
			case BoxToken:
				g.board[y][x] = EmptyToken
				g.actors = append(g.actors, NewBox(x, y))
			case SwitchToken:
				g.board[y][x] = EmptyToken
				g.actors = append(g.actors, NewSwitch(x, y))
			case ClosedDoorToken:
				g.board[y][x] = EmptyToken
				g.actors = append(g.actors, NewDoor(x, y))
			case SpikeUpToken, SpikeDownToken:
				g.board[y][x] = EmptyToken
				g.actors = append(g.actors, NewSpike(x, y, Token(c) == SpikeUpToken))
			default:
				return fmt.Errorf("invalid token: %c", c)
			}
		}
	}

	return nil
}

type StateList map[Actor]*StateChangeNode

type StateChangeNode struct {
	Actor       Actor
	Change      StateChange
	ParentActor Actor
	Parent      *StateChangeNode
	Depth       int
}

func (s *StateChangeNode) String() string {
	var sb strings.Builder
	if s.Parent != nil {
		sb.WriteString(s.Parent.String())
	} else {
		sb.WriteString(fmt.Sprintf("%v ", s.Actor.GetPosition()))
	}
	sb.WriteString(fmt.Sprintf("%v -> %v ", s.Actor, s.Change))
	return sb.String()
}

func (s *StateChangeNode) Equals(other *StateChangeNode) bool {
	if other == nil {
		return false
	}
	if s.Parent != other.Parent {
		return false
	}
	return s.Change.Equals(other.Change)
}

func getAffectingStates(states StateList, actor Actor) map[Actor]StateChange {
	affectingStates := make(map[Actor]StateChange, 0)
	var myChange *StateChangeNode
	myChange = states[actor]

	for a, node := range states {
		if a == actor {
			continue
		}

		pos := actor.GetPosition()
		// something moving onto us
		if node.Change.Move.Equals(pos) {
			affectingStates[a] = node.Change
		}

		// something moving off of us
		if node.Change.From.Equals(pos) {
			affectingStates[a] = node.Change
		}

		// something we're moving to
		if myChange != nil && myChange.Change.Move.Equals(a.GetPosition()) {
			affectingStates[a] = node.Change
		}

		// something updating us
		for _, update := range node.Change.Updates {
			if update == actor {
				affectingStates[a] = node.Change
				break
			}
		}

		// something we're updating
		if myChange != nil {
			for _, update := range myChange.Change.Updates {
				if update == a {
					affectingStates[a] = node.Change
					break
				}
			}
		}
	}
	return affectingStates
}

func (g *Game) getNextStates(states StateList, dir Direction) StateList {
	newStates := make(StateList, 0)
	for _, actor := range g.actors {
		affectingStates := getAffectingStates(states, actor)

		if state, parentActor := actor.Transform(g, dir, affectingStates); state != nil {
			state.From = actor.GetPosition()

			// find parent node
			var parent *StateChangeNode
			depth := 0
			if parentActor != nil {
				parent = states[parentActor]
				depth = parent.Depth + 1
			}

			node := &StateChangeNode{
				Actor:       actor,
				Change:      *state,
				ParentActor: parentActor,
				Parent:      parent,
				Depth:       depth,
			}
			newStates[actor] = node
		}
	}
	return newStates
}

func mergeStates(states StateList, newStates StateList) bool {
	changed := false
	for actor, state := range newStates {
		if oldState, ok := states[actor]; ok {
			if !oldState.Equals(state) {
				changed = true
				states[actor] = state // the state is different from the old one
			}
		} else {
			changed = true
			states[actor] = state // this is a new state
		}

	}

	return changed
}

func getLeaves(states StateList) StateList {
	leaves := make(StateList, 0)
	for a, state := range states {
		// check if anybody has this state as a parent
		found := false
		for _, otherState := range states {
			if otherState.Parent == state {
				found = true
				break
			}
		}

		if !found {
			leaves[a] = state
		}
	}
	return leaves
}

func hasLeaves(states StateList) bool {
	leaves := getLeaves(states)
	return len(leaves) > 0
}

type StateChangeLeaf struct {
	Actor Actor
	Node  *StateChangeNode
}

func popLeaves(states StateList) []StateChangeLeaf {
	leaves := getLeaves(states)
	for actor := range leaves {
		delete(states, actor)
	}

	sortedLeaves := make([]StateChangeLeaf, 0)
	for actor, node := range leaves {
		sortedLeaves = append(sortedLeaves, StateChangeLeaf{actor, node})
	}

	sort.Slice(sortedLeaves, func(i, j int) bool {
		return sortedLeaves[i].Node.Depth > sortedLeaves[j].Node.Depth
	})

	return sortedLeaves
}

func (g *Game) Move(dir Direction) {
	states := make(StateList, 0)
	step := 1
	changed := true
	for changed {
		if step > 10 {
			panic("too many steps")
		}

		g.Println("calculating new states step: ", step)
		step++
		changed = mergeStates(states, g.getNextStates(states, dir))
		if changed {
			leaves := getLeaves(states)
			for _, node := range leaves {
				g.Println(node)
			}
		} else {
			g.Println("no new states")
		}

		g.Println("----------")
	}

	step = 0
	for hasLeaves(states) {
		g.Println("apply states step: ", step)
		step++
		leaves := popLeaves(states)
		for _, leaf := range leaves {
			change := leaf.Node.Change
			actor := leaf.Actor
			g.Println("apply state:", change, " to actor:", actor)
			actor.Apply(g, change)
		}
		g.Println("----------")
	}

	for _, actor := range g.actors {
		actor.Tick(g)
	}

	for _, actor := range g.killQueue {
		g.RemoveActor(actor)
	}
	g.killQueue = make([]Actor, 0)
}
