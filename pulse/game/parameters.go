package game

type Parameter struct {
	add     func()
	sub     func()
	display func() string
}

func Param(display func() string, add func(), sub func()) *Parameter {
	return &Parameter{
		add:     add,
		sub:     sub,
		display: display,
	}
}

type Parameters struct {
	parameters []*Parameter
	cursor     int
}

func NewParameters(p ...*Parameter) *Parameters {
	if len(p) == 0 {
		panic("Parameters must not be empty")
	}

	return &Parameters{
		parameters: p,
		cursor:     0,
	}
}

func (p *Parameters) Add() {
	p.parameters[p.cursor].add()
}

func (p *Parameters) Sub() {
	p.parameters[p.cursor].sub()
}

func (p *Parameters) Display() string {
	return p.parameters[p.cursor].display()
}

func (p *Parameters) Next() {
	p.cursor++
	if p.cursor >= len(p.parameters) {
		p.cursor = 0
	}
}

func (p *Parameters) Prev() {
	p.cursor--
	if p.cursor < 0 {
		p.cursor = len(p.parameters) - 1
	}
}
