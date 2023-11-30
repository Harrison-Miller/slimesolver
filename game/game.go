package game

import (
	"fmt"
	"strings"
)

type Token rune

const (
	WallToken       Token = '#'
	EmptyToken      Token = '.'
	PitToken        Token = 'O'
	SwitchToken     Token = 'X'
	ClosedDoorToken Token = 'D'
	OpenDoorToken   Token = '_'
	BoxToken        Token = 'B'
	SlimeToken      Token = '@'
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

func (g *Game) GetActors(x, y int) []Actor {
	l := make([]Actor, 0)
	for _, entity := range g.actors {
		v := entity.GetPosition()
		if v.X == x && v.Y == y {
			l = append(l, entity)
		}
	}
	return l
}

func (g *Game) GetActorsWithToken(token Token) []Actor {
	l := make([]Actor, 0)
	for _, actor := range g.actors {
		if actor.Token() == token {
			l = append(l, actor)
		}
	}
	return l
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
			actors := g.GetActors(x, y)
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
			case SlimeToken:
				g.board[y][x] = EmptyToken
				g.actors = append(g.actors, NewSlime(x, y))
			case BoxToken:
				g.board[y][x] = EmptyToken
				g.actors = append(g.actors, NewBox(x, y))
			case SwitchToken:
				g.board[y][x] = EmptyToken
				g.actors = append(g.actors, NewSwitch(x, y))
			case ClosedDoorToken:
				g.board[y][x] = EmptyToken
				g.actors = append(g.actors, NewDoor(x, y))
			default:
				return fmt.Errorf("invalid token: %c", c)
			}
		}
	}

	return nil
}

func (g *Game) extendGraph(graph *Graph, previous *ActorNode, current *LocationNode) {
	for _, actorNode := range current.Actors {
		if actorNode.Edges != nil && len(actorNode.Edges) > 0 {
			for _, edge := range actorNode.Edges {
				g.extendGraph(graph, actorNode, edge)
			}
		} else if previous != nil { // extend leaf nodes
			// TODO: since we're replacing the edges, the order of this matters, so we'll need to choose based on depth
			dir := directionBetween(previous.Actor.GetPosition(), actorNode.Actor.GetPosition())
			newEdges := actorNode.Actor.CalculateEdges(g, dir, previous.Actor)
			if len(newEdges) != 0 {
				graph.UpdateActorEdges(actorNode, newEdges)
				for _, edge := range actorNode.Edges {
					g.extendGraph(graph, actorNode, edge)
				}
			}
		}
	}
}

func (g *Game) Move(dir Direction) {
	// build gameGraph
	gameGraph := NewGraph()
	for _, entity := range g.actors {
		edges := entity.CalculateEdges(g, dir, nil)
		gameGraph.AddActorNode(NewActorNode(entity, edges))
	}
	gameGraph.Compute()
	fmt.Printf("built gameGraph for move (%s):\n", dirString(dir))
	fmt.Println(gameGraph)

	// extend gameGraph
	roots := gameGraph.GetRoots()
	for _, root := range roots {
		g.extendGraph(gameGraph, nil, root)
	}
	fmt.Println("extended gameGraph:")
	fmt.Println(gameGraph)

	for gameGraph.HasLeaves() {
		leaves := gameGraph.PopLeaves()
		for _, leaf := range leaves {
			if leaf.Actor != nil {
				leaf.Actor.ApplyEdges(g, leaf.EdgePositions)
			}
		}
	}

	// TODO: this probably needs to be sorted somehow
	for _, actor := range g.actors {
		actor.ResolveState(g)
	}
	fmt.Println("--------------------")
}
