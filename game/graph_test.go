package game

import (
	"fmt"
	"slimesolver/game/math"
	"testing"
)

type TestActor struct {
	X, Y int
}

func (t *TestActor) Token() Token {
	return EmptyToken
}

func (t *TestActor) GetPosition() math.Vector2 {
	return math.Vector2{t.X, t.Y}
}

func (t *TestActor) CalculateEdges(g *Game, dir Direction, a Actor) []math.Vector2 {
	return []math.Vector2{}
}

func (t *TestActor) ApplyEdges(g *Game, edges []math.Vector2) {

}

func (t *TestActor) ResolveState(g *Game) {

}

func (t *TestActor) Solid() bool {
	return false
}

func singleEdge(x, y int) []math.Vector2 {
	return []math.Vector2{{x, y}}
}

func TestGraph(t *testing.T) {
	t.Run("two separate nodes", func(t *testing.T) {
		graph := NewGraph()
		graph.AddActorNode(NewActorNode(&TestActor{0, 0}, singleEdge(2, 0)))
		graph.AddActorNode(NewActorNode(&TestActor{1, 0}, singleEdge(3, 0)))
		graph.Compute()
		fmt.Println(graph)

		roots := graph.GetRoots()
		if len(roots) != 2 {
			t.Errorf("expected 2 roots, got %v", len(roots))
			return
		}

		node1 := roots[0]
		if node1.String() != "(0, 0) -> . -> (2, 0)" {
			t.Errorf("expected (0, 0) -> . -> (2, 0), got %v", node1.String())
		}

		node2 := roots[1]
		if node2.String() != "(1, 0) -> . -> (3, 0)" {
			t.Errorf("expected (1, 0) -> . -> (3, 0), got %v", node2.String())
		}

		// a leaf is a node that does not effect any others
		leaves := graph.GetLeaves()
		if len(leaves) != 2 {
			t.Errorf("expected 2 leaves, got %v", len(leaves))
			return
		}

		if leaves[0].String() != ". -> (2, 0)" {
			t.Errorf("expected . -> (2, 0), got %v", leaves[0].String())
		}

		if leaves[1].String() != ". -> (3, 0)" {
			t.Errorf("expected . -> (3, 0), got %v", leaves[1].String())
		}
	})

	t.Run("2 nodes pushing the same direction", func(t *testing.T) {
		graph := NewGraph()
		graph.AddActorNode(NewActorNode(&TestActor{0, 0}, singleEdge(1, 0)))
		graph.AddActorNode(NewActorNode(&TestActor{1, 0}, singleEdge(2, 0)))
		graph.Compute()
		fmt.Println(graph)

		roots := graph.GetRoots()
		if len(roots) != 1 {
			t.Errorf("expected 1 roots, got %v", len(roots))
			return
		}

		node := roots[0]
		if node.String() != "(0, 0) -> . -> (1, 0) -> . -> (2, 0)" {
			t.Errorf("expected (0, 0) -> . -> (1, 0) -> . -> (2, 0), got %v", node.String())
		}

		leaves := graph.GetLeaves()
		if len(leaves) != 1 {
			t.Errorf("expected 1 leaves, got %v", len(leaves))
			return
		}

		if leaves[0].String() != ". -> (2, 0)" {
			t.Errorf("expected . -> (2, 0), got %v", leaves[0].String())
		}
	})

	t.Run("adding in reverse order", func(t *testing.T) {
		graph := NewGraph()
		graph.AddActorNode(NewActorNode(&TestActor{0, 0}, singleEdge(1, 0)))
		graph.AddActorNode(NewActorNode(&TestActor{1, 0}, singleEdge(2, 0)))
		graph.Compute()
		fmt.Println(graph)

		roots := graph.GetRoots()
		if len(roots) != 1 {
			t.Errorf("expected 1 roots, got %v", len(roots))
			return
		}

		node := roots[0]
		if node.String() != "(0, 0) -> . -> (1, 0) -> . -> (2, 0)" {
			t.Errorf("expected (0, 0) -> . -> (1, 0) -> . -> (2, 0), got %v", node.String())
		}

		leaves := graph.GetLeaves()
		if len(leaves) != 1 {
			t.Errorf("expected 1 leaves, got %v", len(leaves))
			return
		}

		if leaves[0].String() != ". -> (2, 0)" {
			t.Errorf("expected . -> (2, 0), got %v", leaves[0].String())
		}
	})

	t.Run("2 nodes, connect to the same place", func(t *testing.T) {
		graph := NewGraph()
		graph.AddActorNode(NewActorNode(&TestActor{0, 0}, singleEdge(2, 0)))
		graph.AddActorNode(NewActorNode(&TestActor{1, 0}, singleEdge(2, 0)))
		graph.AddActorNode(NewActorNode(&TestActor{2, 0}, nil))
		graph.Compute()
		fmt.Println(graph)

		roots := graph.GetRoots()
		if len(roots) != 2 {
			t.Errorf("expected 2 roots, got %v", len(roots))
			return
		}

		if roots[0].String() != "(0, 0) -> . -> (2, 0) -> ." {
			t.Errorf("expected (0, 0) -> . -> (2, 0) -> ., got %v", roots[0].String())
		}

		if roots[1].String() != "(1, 0) -> . -> (2, 0) -> ." {
			t.Errorf("expected (1, 0) -> . -> (2, 0) -> ., got %v", roots[1].String())
		}

		leaves := graph.GetLeaves()
		if len(leaves) != 1 {
			t.Errorf("expected 1 leaf, got %v", len(leaves))
			return
		}

		if leaves[0].String() != "." {
			t.Errorf("expected ., got %v", leaves[0].String())
		}
	})

	t.Run("going in reverse order", func(t *testing.T) {
		graph := NewGraph()
		graph.AddActorNode(NewActorNode(&TestActor{1, 0}, singleEdge(0, 0)))
		graph.AddActorNode(NewActorNode(&TestActor{2, 0}, singleEdge(1, 0)))
		graph.Compute()
		fmt.Println(graph)

		roots := graph.GetRoots()
		if len(roots) != 1 {
			t.Errorf("expected 1 roots, got %v", len(roots))
			return
		}

		node := roots[0]
		if node.String() != "(2, 0) -> . -> (1, 0) -> . -> (0, 0)" {
			t.Errorf("expected (2, 0) -> . -> (1, 0) -> . -> (0, 0), got %v", node.String())
		}

		leaves := graph.GetLeaves()
		if len(leaves) != 1 {
			t.Errorf("expected 1 leaves, got %v", len(leaves))
			return
		}

		if leaves[0].String() != ". -> (0, 0)" {
			t.Errorf("expected . -> (0, 0), got %v", leaves[0].String())
		}
	})
}

func TestGraphPopLeaves(t *testing.T) {
	graph := NewGraph()
	graph.AddActorNode(NewActorNode(&TestActor{0, 0}, singleEdge(1, 0)))
	graph.AddActorNode(NewActorNode(&TestActor{1, 0}, singleEdge(2, 0)))
	graph.Compute()
	fmt.Println(graph)

	leaves := graph.PopLeaves()
	if len(leaves) != 1 {
		t.Errorf("expected 1 leaf, got %v", len(leaves))
		return
	}

	if leaves[0].String() != ". -> (2, 0)" {
		t.Errorf("expected . -> (2, 0), got %v", leaves[0].String())
	}

	fmt.Println(graph)
	leaves = graph.PopLeaves()
	if len(leaves) != 1 {
		t.Errorf("expected 1 leaf, got %v", len(leaves))
		return
	}

	if leaves[0].String() != ". -> (1, 0)" {
		t.Errorf("expected . -> (1, 0), got %v", leaves[0].String())
	}
}
