package xmlconv
import (
	"github.com/Andrew-M-C/go-tools/str"
)

func (x *Item) GetChild(name string, names ...string) (c *Item, exist bool) {
	cs, exist := x.child[name]
	if false == exist {
		return nil, false
	}

	switch len(names) {
	case 0:
		return cs[0], true
	case 1:
		return cs[0].GetChild(names[0])
	default:
		return cs[0].GetChild(names[0], names[1:]...)
	}
}


func (x *Item) GetChildString(name string, names ...string) (s string, exist bool) {
	c, exist := x.GetChild(name, names...)
	if false == exist {
		return "", false
	} else {
		return c.String(), true
	}
}


func (x *Item) GetChildBytes(name string, names ...string) (b []byte, exist bool) {
	c, exist := x.GetChild(name, names...)
	if false == exist {
		return nil, false
	} else {
		return c.Bytes(), true
	}
}


func (x *Item) GetChildren(name string) (c []*Item, exist bool) {
	c, exist = x.child[name]
	if false == exist {
		return nil, false
	}
	return c, true
}


func (x *Item) GetChildAtIndex(name string, i int) (c *Item, exist bool) {
	cs, exist := x.child[name]
	if false == exist {
		return nil, false
	}
	if i >= len(cs) {
		return nil, false
	}
	return cs[i], true
}


func (x *Item) SetChild(child *Item, name string, names ...string) *Item {
	if nil == child || str.Empty(name) {
		return nil
	}

	if 0 == len(names) {
		child.name = name
		children := make([]*Item, 0, 1)
		children = append(children, child)
		x.child[name] = children
		return child

	} else {
		var c *Item
		cs, exist := x.child[name]
		if false == exist {
			c = NewItem(name)
			cs = append(make([]*Item, 0, 1), c)
			x.child[name] = cs
		} else {
			c = cs[0]
		}

		if len(names) > 1 {
			return c.SetChild(child, names[0], names[1:]...)
		} else {
			return c.SetChild(child, names[0])
		}
	}
}


func (x *Item) AddChild(child *Item, name string) *Item {
	if nil == child || str.Empty(name) {
		return nil
	}

	child.name = name
	children, exist := x.child[name]
	if false == exist {
		children = make([]*Item, 0, 1)
	}
	children = append(children, child)
	x.child[name] = children
	return child
}


func (x *Item) SetEmptyChild(name string, names ...string) *Item {
	if str.Empty(name) {
		return nil
	}
	return x.SetChild(NewItem(name), name, names...)
}


func (x *Item) SetBytesChild(b []byte, name string, names ...string) *Item {
	if str.Empty(name) {
		return nil
	}

	c := NewItem(name)
	c.data = b
	return x.SetChild(c, name, names...)
}


func (x *Item) SetStringChild(s string, name string, names ...string) *Item {
	if str.Empty(name) {
		return nil
	}

	c := NewItem(name)
	c.data = []byte(s)
	c.dataString = &s
	return x.SetChild(c, name, names...)
}
