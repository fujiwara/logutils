package logutils

import (
	"bytes"

	"github.com/fatih/color"
)

func Color(attr ...color.Attribute) ModifierFunc {
	c := color.New(attr...)
	buf := &bytes.Buffer{}
	return func(b []byte) []byte {
		buf.Reset()
		c.Fprint(buf, string(b))
		return buf.Bytes()
	}
}
