package ui

type Component interface {
	GetWidth() int
	GetHeight() int
}

const (
	// PositionCenter components are centered in their parent container.
	PositionCenter = iota

	// PositionTop components are aligned to the top of their parent container.
	PositionTop

	// PositionBottom components are aligned to the bottom of their parent container.
	PositionBottom

	// PositionLeft components are aligned to the left of their parent container.
	PositionLeft

	// PositionRight components are aligned to the right of their parent container.
	PositionRight

	// PositionAbsolute positioning not handled by the layout system.
	PositionAbsolute
)

const (
	// SizeNone sizing not handled by the layout system.
	SizeNone = iota

	// SizeFill components fill their parent container.
	SizeFill

	// SizeRelative components size relative to their parent container.
	SizeRelative
)

type Node struct {
	Component
	Children []*Node
	Margin   [4]int // [top, right, bottom, left]
	Padding  [4]int // [top, right, bottom, left]
	Position [2]int // [x, y]
	Sizing   [2]int // [width, height]
}

func NewNode(c Component) *Node {
	return &Node{
		Children:  make([]*Node, 0),
		Margin:    [4]int{0, 0, 0, 0},
		Padding:   [4]int{0, 0, 0, 0},
		Position:  [2]int{PositionAbsolute, PositionAbsolute},
		Sizing:    [2]int{SizeNone, SizeNone},
		Component: c,
	}
}

type Layout struct {
	Root *Node
}
