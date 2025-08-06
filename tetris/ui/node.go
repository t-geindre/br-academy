package ui

type Component interface {
	SetSize(width, height int)
	SetPosition(x, y int)
}

const (
	OrientationHorizontal = iota
	OrientationVertical

	PositionStart
	PositionEnd
	PositionCenter

	UnitPixel
	UnitPercentage
)

type Node struct {
	Component
	Children []*Node
	Parent   *Node

	Orientation uint8
	Position    uint8

	Padding     [2]int
	PaddingUnit uint8

	Margin     [2]int
	MarginUnit uint8

	MinSize int
	MaxSize int

	W, H, X, Y int
}

func NewNode(component Component) *Node {
	return &Node{
		Component:   component,
		Children:    []*Node{},
		Orientation: OrientationVertical,
		Position:    PositionCenter,
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
