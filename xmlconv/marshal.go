package xmlconv
import (
	"github.com/Andrew-M-C/go-tools/str"
	"strings"
	"bytes"
)

type Option struct {
	Indent	string
}

var defaultOpt = Option{
	Indent: "",
}


func (self *Item) MarshalBytes(opt ...Option) ([]byte, error) {
	indent := ""
	if len(opt) > 0 {
		indent = opt[0].Indent
	} else {
		indent = defaultOpt.Indent
	}
	buff := bytes.Buffer{}
	self.toBuffer(&buff, indent, 0)
	return buff.Bytes(), nil
}


func (self *Item) MarshalString(opt ...Option) (string, error) {
	indent := ""
	if len(opt) > 0 {
		indent = opt[0].Indent
	} else {
		indent = defaultOpt.Indent
	}
	buff := bytes.Buffer{}
	self.toBuffer(&buff, indent, 0)
	return buff.String(), nil
}


func (self *Item) Marshal(opt ...Option) (string, error) {
	indent := ""
	if len(opt) > 0 {
		indent = opt[0].Indent
	} else {
		indent = defaultOpt.Indent
	}
	buff := bytes.Buffer{}
	self.toBuffer(&buff, indent, 0)
	return buff.String(), nil
}


func (self *Item) toBuffer(buff *bytes.Buffer, indent string, depth int) {
	prefix := ""
	if str.Valid(indent) {
		prefix = "\n" + strings.Repeat(indent, depth)
	}
	buff.WriteString(prefix)
	buff.WriteRune('<')
	buff.WriteString(self.name)

	for k, v := range self.attrs {
		buff.WriteRune(' ')
		buff.WriteString(k)
		buff.WriteRune('=')
		buff.WriteString(v)
	}
	buff.WriteRune('>')

	s := self.String()
	if str.Valid(s) {
		if strings.ContainsAny(s, "<>&\"'") {
			buff.WriteString("<![CDATA[")
			buff.WriteString(s)
			buff.WriteString("]]>")
		} else {
			buff.WriteString(s)
		}
	}

	if len(self.child) > 0 {
		for _, c := range self.child {
			c.toBuffer(buff, indent, depth + 1)
		}
		buff.WriteString(prefix)
	}

	buff.WriteString("</")
	buff.WriteString(self.name)
	buff.WriteRune('>')
}
