/*
 * reference:
 * - [定制实现 Go 中的 XML Unmarshal - 基础篇](https://github.com/imjoey/blog/issues/19)
 */

package xmlconv
import (
	// "github.com/Andrew-M-C/go-tools/str"
	// "github.com/Andrew-M-C/go-tools/log"
	"encoding/xml"
	"bytes"
	"io"
)

type Item struct {
	name		string
	data		[]byte
	dataString	*string
	attrs		map[string]string
	child		map[string]*Item
}

func (x *Item) Name() string {
	return x.name
}

func (x *Item) Bytes() []byte {
	return x.data
}

func (x *Item) SetData(b []byte) {
	x.data = b
	x.dataString = nil
	return
}

func (x *Item) SetString(s string) {
	x.data = []byte(s)
	x.dataString = nil
	return
}

func (x *Item) String() string {
	if nil == x.dataString {
		s := string(x.data)
		x.dataString = &s
	}
	return *(x.dataString)
}

func (x *Item) Attrs() map[string]string {
	return x.attrs
}

func (x *Item) GetAttr(a string, defaultValue ...string) (string, bool) {
	ret, exist := x.attrs[a]
	if false == exist && len(defaultValue) > 0 {
		ret = defaultValue[0]
	}
	return ret, exist
}

func (x *Item) SetAttr(n, v string) {
	x.attrs[n] = v
	return
}

func (x *Item) Children() map[string]*Item {
	return x.child
}


func NewItem(name string) *Item {
	// log.Debug("NewItem %s", name)
	ret := Item{
		name: name,
		attrs: make(map[string]string),
		child: make(map[string]*Item),
		data: []byte(""),
	}
	return &ret
}


func NewFromString(s string) (*Item, error) {
	return NewFromBytes([]byte(s))
}


func NewFromBytes(b []byte) (*Item, error) {
	if 0 == len(b) {
		return nil, ParaError
	}

	// log.Debug("input string: %s", string(b))
	decoder := xml.NewDecoder(bytes.NewReader(b))
	stk := newStack()
	var curr *Item
	var root *Item

	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				// log.Debug("Parse XML finished!")
				return root, nil
			} else {
				// log.Debug("Failed to Parse XML with the error of %v", err)
				return nil, err
			}
			break
		}
		t = xml.CopyToken(t)

		switch t := t.(type) {
		case xml.StartElement:
			// log.Debug("xml.StartElement")
			name := t.Name.Local
			item := NewItem(name)
			attr := item.attrs
			// log.Debug("name: %s", name)

			for _, a := range t.Attr {
				// log.Debug("attr: %s - %s", a.Name.Local, a.Value)
				attr[a.Name.Local] = a.Value
			}
			if curr != nil {
				curr.child[name] = item
				stk.Push(curr)
			} else {
				root = item
			}
			curr = item

		case xml.EndElement:
			// log.Debug("xml.EndElement")
			curr = stk.Pop()

		case xml.CharData:
			b := []byte(t)
			b = bytes.Trim(b, "\r\n\t ")
			if false == bytesEmpty(b) {
				// log.Debug("xml.CharData: '%s'", string(b))
				if curr != nil {
					curr.data = b
				}
			}

		case xml.Comment:
			// ignore
		}
	}

	return nil, FormatError
}

func bytesEmpty(b []byte) bool {
	for _, _ = range b {
		return false
	}
	return true
}
