package pool

type Pool struct {
	items []any
}

func NewPool() *Pool {
	return &Pool{
		items: make([]any, 0),
	}
}

func (p *Pool) Add(item any) {
	p.items = append(p.items, item)
}

func (p *Pool) Remove(item any) {
	for i, v := range p.items {
		if v == item {
			p.items = append(p.items[:i], p.items[i+1:]...)
			return
		}
	}
}
