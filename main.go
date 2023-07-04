package main

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y int
}

type Node struct {
	box            Box
	Point          *Point
	NE, NW, SE, SW *Node
}

type QuadTree struct {
	Root *Node
}

type Branch []*Node

// Insert a new point into the QuadTree
func (qt *QuadTree) Insert(p *Point) {
	// Insertion logic here
}

// Search for a point in the QuadTree
func (qt *QuadTree) Search(p *Point) *Node {
	// Search logic here
	return nil
}

// Retrieve the nearest point to a given point in the QuadTree
func (qt *QuadTree) Nearest(p *Point) *Point {

	branch := qt.Root.BranchFromPoint(p)
	if len(branch) == 0 {
		return nil
	}
	nearest, _ := Range(branch)
	return nearest

}

// BranchFromPoint retrieves a path where each path node contains recursively the point.
// Is important to note that this function is executed from a given node.
func (n *Node) BranchFromPoint(p *Point) Branch {
	if !n.box.ContainsPoint(p) {
		return []*Node{}
	}
	if n.isLeaf() {
		return []*Node{n}
	}
	r := make([]*Node, 0)
	r = append(r, n)

	for _, child := range n.getChildren() {
		if child != nil && child.box.ContainsPoint(p) {
			r = append(r, child.BranchFromPoint(p)...)
		}
	}
	return r
}

func Range(branch Branch) (*Point, float64) {
	nearest, distance := &Point{}, math.MaxFloat64
	for j := len(branch) - 1; j >= 0; j-- {
		current := branch[j]
		p, d := current.NearestPoint(distance, branch.nextNode(j))
		if d < distance {
			nearest, distance = p, d
		}
	}
	return nearest, distance
}

// NearestPoint given a node and a cap determine the nearest point that intersects with the cap radius.
func (n Node) NearestPoint(cap *Cap, avoid *Node) (*Point, float64) {
	nearest, distance := &Point{}, math.MaxFloat64
	if n.isLeaf() {
		d := cap.Center.distance(n.Point)
		if d <= cap.Radius {
			nearest, distance = n.Point, d
			cap.Expand(d)
		}

	} else {
		for _, c := range n.getChildren() {
			if c != nil && c != avoid && cap.IntersectsRect(c.box.rectangle) {
				p, d := c.NearestPoint(cap, avoid)
				if d < distance {
					nearest, distance = p, d
				}
			}
		}
	}
	return nearest, distance
}

func (n Node) isLeaf() bool {
	return n.NE == nil && n.NW == nil && n.SE == nil && n.SW == nil
}

func (n Node) getChildren() []*Node {
	return []*Node{n.NE, n.NW, n.SE, n.SW}
}

func (p1 *Point) distance(p2 *Point) float64 {
	return math.Sqrt(math.Pow(float64(p1.X-p2.X), 2) + math.Pow(float64(p1.Y-p2.Y), 2))
}

// lastNode retrieve the lastNode element of the branch.
func (b Branch) lastNode() *Node {
	if len(b) == 1 {
		return b[0]
	}
	return b[len(b)-1]
}

// nextNode retrieve the next node element of the branch.
func (b Branch) nextNode(i int) *Node {
	lastIndex := len(b) - 1
	if i+1 >= lastIndex {
		return b[lastIndex]
	}
	if i+1 <= 0 {
		return b[0]
	}
	return b[i+1]
}

func main() {
	qt := &QuadTree{}

	point := &Point{X: 5, Y: 5}
	qt.Insert(point)

	fmt.Println(qt.Search(point))
	fmt.Println(qt.Nearest(point))
}
