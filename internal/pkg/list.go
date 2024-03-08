package pkg

import (
	"container/list"
)

func NextN(l *list.List, p *list.Element, n int) *list.Element {
	if n > 0 {
		for i := 0; i < n; i++ {
			p = NextForward(l, p)
		}
		return p
	}

	n = -n
	for i := 0; i < n; i++ {
		p = NextBackward(l, p)
	}
	return p
}

func NextForward(l *list.List, p *list.Element) *list.Element {
	if p.Next() != nil {
		return p.Next()
	}

	return l.Front()
}

func NextBackward(l *list.List, p *list.Element) *list.Element {
	if p.Prev() != nil {
		return p.Prev()
	}

	return l.Back()
}
