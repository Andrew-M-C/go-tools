package xmlconv
import ()

type stack struct {
	L	[]*Item
}

func newStack() *stack {
	ret := stack{
		L: make([]*Item, 0, 0),
	}
	return &ret
}


func (s *stack) Push(i *Item) {
	if i != nil {
		s.L = append(s.L, i)
	}
}


func (s *stack) Pop() *Item {
	l := len(s.L)
	if l <= 0 {
		return nil
	} else {
		ret := s.L[l-1]
		s.L = s.L[0:l-1]
		return ret
	}
}
