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
	if depth > 0 {
		buff.WriteString(prefix)
	}
	buff.WriteRune('<')
	buff.WriteString(self.name)

	for k, v := range self.attrs {
		buff.WriteRune(' ')
		buff.WriteString(k)
		buff.WriteString("=\"")
		writeAttrToBuff(v, buff)
		buff.WriteRune('"')
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


func writeAttrToBuff(v string, buff *bytes.Buffer) {
	for _, c := range v {
		switch c {
		case '&':
			buff.WriteString("&amp;")
		case '<':
			buff.WriteString("&lt;")
		case '>':
			buff.WriteString("&gt;")
		case '"':
			buff.WriteString("&quot;")
		case '\'':
			buff.WriteString("&apos;")
		default:
			buff.WriteRune(c)
		}
	}
	return
}
