package game

import (
	"slimesolver/game/math"
	"strings"
)

type LocationNode struct {
	Position math.Vector2
	Actors   []*ActorNode
}

func (n *LocationNode) String() string {
	var sb strings.Builder
	sb.WriteString(n.Position.String())

	for i, actor := range n.Actors {
		if i == 0 {
			sb.WriteString(" -> ")
		} else {
			sb.WriteString("\t -> ")
		}
		sb.WriteString(actor.String())
	}
	return sb.String()
}

type ActorNode struct {
	Position      math.Vector2
	Actor         Actor
	EdgePositions []math.Vector2
	Edges         []*LocationNode
	Depth         int
}

func NewActorNode(a Actor, edges []math.Vector2) *ActorNode {
	return &ActorNode{
		Position:      a.GetPosition(),
		Actor:         a,
		EdgePositions: edges,
		Edges:         make([]*LocationNode, 0),
	}
}

func (n *ActorNode) String() string {
	var sb strings.Builder
	sb.WriteString(string(n.Actor.Token()))
	if n.Edges != nil && len(n.Edges) > 0 {
		for i, edge := range n.Edges {
			if i == 0 {
				sb.WriteString(" -> ")
			} else {
				sb.WriteString("\t -> ")
			}
			sb.WriteString(edge.String())
		}
	} else {
		for i, edge := range n.EdgePositions {
			if i == 0 {
				sb.WriteString(" -> ")
			} else {
				sb.WriteString("\t -> ")
			}
			sb.WriteString(edge.String())
		}
	}

	return sb.String()
}

type Graph struct {
	nodes []*LocationNode
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make([]*LocationNode, 0),
	}
}

func (g *Graph) String() string {
	roots := g.GetRoots()
	var sb strings.Builder
	for i, root := range roots {
		if i != 0 {
			sb.WriteString("\n---\n")
		}
		sb.WriteString(root.String())

	}
	return sb.String()
}

func (g *Graph) AddActorNode(node *ActorNode) {
	if node == nil {
		return
	}

	locationNode := g.GetLocationNode(node.Position)
	// need to create a new location node
	if locationNode == nil {
		locationNode = &LocationNode{
			Position: node.Position,
			Actors:   []*ActorNode{node},
		}
		g.nodes = append(g.nodes, locationNode)
		return
	}

	// add the actors to the current location node
	locationNode.Actors = append(locationNode.Actors, node)
}

func (g *Graph) UpdateActorEdges(node *ActorNode, edges []math.Vector2) {
	if node == nil {
		return
	}

	node.EdgePositions = edges
	node.Edges = make([]*LocationNode, 0)

	for _, edgePosition := range edges {
		edgeNode := g.GetLocationNode(edgePosition)
		if edgeNode != nil {
			node.Edges = append(node.Edges, edgeNode)
		}
	}

	// TODO: recalculate depth
}

func (g *Graph) RemoveLocationNode(node *LocationNode) {
	if node == nil {
		return
	}

	// remove from list
	for i, otherNode := range g.nodes {
		if otherNode == node {
			g.nodes = append(g.nodes[:i], g.nodes[i+1:]...)
			break
		}
	}

	// remove from any edges
	for _, otherNode := range g.nodes {
		for _, actor := range otherNode.Actors {
			for j, edge := range actor.Edges {
				if edge == node {
					actor.Edges = append(actor.Edges[:j], actor.Edges[j+1:]...)
					break
				}
			}
		}
	}
}

func (g *Graph) RemoveActorNode(node *ActorNode) {
	locationNode := g.GetLocationNode(node.Position)
	if locationNode == nil {
		return
	}
	for i, actor := range locationNode.Actors {
		if actor == node {
			locationNode.Actors = append(locationNode.Actors[:i], locationNode.Actors[i+1:]...)
			if len(locationNode.Actors) == 0 {
				g.RemoveLocationNode(locationNode)
			}
			break
		}
	}
}

func (g *Graph) GetLocationNode(position math.Vector2) *LocationNode {
	for _, node := range g.nodes {
		if node.Position == position {
			return node
		}
	}
	return nil
}

func (g *Graph) Compute() {
	// clear all actor edges
	for _, node := range g.nodes {
		for _, actor := range node.Actors {
			actor.Edges = make([]*LocationNode, 0)
			actor.Depth = 0
		}
	}

	// compute actor edges
	for _, node := range g.nodes {
		for _, actor := range node.Actors {
			for _, edgePosition := range actor.EdgePositions {
				edgeNode := g.GetLocationNode(edgePosition)
				if edgeNode != nil {
					actor.Edges = append(actor.Edges, edgeNode)
				}
			}
		}
	}

	// compute actor depth
	//roots := g.GetRoots()
	//for _, root := range roots {
	//	g.computeDepth(root, 0)
	//}
}

func (g *Graph) computeDepth(node *LocationNode, depth int) {
	if node == nil {
		return
	}

	for _, actor := range node.Actors {
		if actor.Depth < depth {
			actor.Depth = depth
		}
	}
	for _, edge := range node.Actors {
		locationNode := g.GetLocationNode(edge.Position)
		g.computeDepth(locationNode, depth+1)
	}
}

func (g *Graph) GetRoots() []*LocationNode {
	// a root has no edges to it
	roots := make([]*LocationNode, 0)
	for _, node := range g.nodes {
		found := false
		for _, otherNode := range g.nodes {
			for _, actor := range otherNode.Actors {
				for _, edge := range actor.Edges {
					if edge == node {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			roots = append(roots, node)
		}
	}

	return roots
}

func (g *Graph) GetLeaves() []*ActorNode {
	// a leaf has no edges from it
	leaves := make([]*ActorNode, 0)
	for _, node := range g.nodes {
		for _, actor := range node.Actors {
			if len(actor.Edges) == 0 {
				leaves = append(leaves, actor)
			}
		}
	}
	return leaves
}

func (g *Graph) HasLeaves() bool {
	leaves := g.GetLeaves()
	return len(leaves) > 0
}

func (g *Graph) PopLeaves() []*ActorNode {
	leaves := g.GetLeaves()
	for _, leaf := range leaves {
		g.RemoveActorNode(leaf)
	}

	// TODO: sort leaves by lowest depth first
	//sort.Slice(leaves, func(i, j int) bool {
	//	return leaves[i].Depth < leaves[j].Depth
	//})

	return leaves
}
