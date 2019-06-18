package xmlconv
import (
	"github.com/Andrew-M-C/go-tools/str"
)

func (x *Item) GetChild(n1 string, names ...string) (c *Item, exist bool) {
	c, exist = x.child[n1]

	l := len(names)
	if 0 == l {
		return c, exist
	}

	for _, n := range names {
		c, exist = c.child[n]
		if false == exist {
			return nil, exist
		}
	}
	return c, exist
}


func (x *Item) SetChild(child *Item, n1 string, names ...string) *Item {
	if nil == child || str.Empty(n1) {
		return nil
	}

	l := len(names)
	if 0 == l {
		child.name = n1
		x.child[n1] = child
		return child
	}

	c, exist := x.child[n1]
	if false == exist {
		c = NewItem(n1)
		x.child[n1] = c
	}
	for i, n := range names {
		if i == l - 1 {
			c.child[n] = child
			child.name = n
		} else {
			new_c, exist := c.child[n]
			if false == exist {
				new_c = NewItem(n)
				c.child[n] = new_c
			}
			c = new_c
		}
	}
	return c
}


func (x *Item) SetEmptyChild(n1 string, names ...string) *Item {
	if str.Empty(n1) {
		return nil
	}

	l := len(names)
	if 0 == l {
		return x.SetChild(NewItem(n1), n1)
	} else {
		n := names[l - 1]
		return x.SetChild(NewItem(n), n1, names ...)
	}
}


func (x *Item) SetChildString(s string, n1 string, names ...string) *Item {
	return x.SetChildBytes([]byte(s), n1, names...)
}


func (x *Item) SetChildBytes(b []byte, n1 string, names ...string) *Item {
	if str.Empty(n1) {
		return nil
	}

	l := len(names)
	c, exist := x.child[n1]
	if false == exist {
		c = NewItem(n1)
		x.child[n1] = c
	}

	if 0 == l {
		c.data = b
		return c
	}

	for i, n := range names {
		new_c, exist := c.child[n]
		if false == exist {
			new_c = NewItem(n)
			c.child[n] = new_c
		}
		c = new_c
		if i == l - 1 {
			c.data = b
		}
	}
	return c
}
