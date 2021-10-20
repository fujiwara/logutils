package logutils

import (
	"bytes"
	"log"
	"testing"
)

var testModifier = func(b []byte) []byte {
	r := []byte("xxx")
	r = append(r, b...)
	return r
}

func TestModifierFilter(t *testing.T) {
	buf := new(bytes.Buffer)
	filter := &LevelFilter{
		Levels: []LogLevel{"DEBUG", "WARN", "ERROR"},
		ModifierFuncs: []ModifierFunc{
			testModifier,
			testModifier,
			nil, // no modifier for ERROR
		},
		MinLevel: "WARN",
		Writer:   buf,
	}
	logger := log.New(filter, "", 0)
	logger.Print("[WARN] foo")
	logger.Println("[ERROR] bar")
	logger.Println("[DEBUG] baz")
	logger.Println("[WARN] buzz")

	result := buf.String()
	expected := "xxx[WARN] foo\n[ERROR] bar\nxxx[WARN] buzz\n"
	if result != expected {
		t.Fatalf("bad: %#v", result)
	}
}
