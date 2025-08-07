package ui

type Component interface {
	SetSize(width, height int)
	SetPosition(x, y int)
}

const (
	OrientationHorizontal = iota
	OrientationVertical

	UnitPixel
	UnitPercentage
)

type Node struct {
	Component
	Children []*Node
	Parent   *Node

	ContentOrientation uint8
	ContentSpacing     int
	ContentSpacingUnit uint8

	Padding     [4]int
	PaddingUnit uint8

	Size int // Preferred size px

	Grow   float64 // Grow factor for the node
	Shrink float64 // Shrink factor for the node

	W, H, X, Y int

	Name string // Optional name for debugging
}

func NewNode(component Component) *Node {
	return &Node{
		Component:          component,
		Children:           []*Node{},
		ContentOrientation: OrientationVertical,
		ContentSpacingUnit: UnitPixel,
		PaddingUnit:        UnitPixel,
	}
}

func (n *Node) Append(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

func (n *Node) Remove(child *Node) {
	for i, c := range n.Children {
		if c == child {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			return
		}
	}
	child.Parent = nil
}
