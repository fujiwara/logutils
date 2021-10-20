package logutils

import (
	"bytes"
	"log"
	"testing"

	"github.com/fatih/color"
)

func TestColoerModifierFilter(t *testing.T) {
	color.NoColor = false

	buf := new(bytes.Buffer)
	filter := &LevelFilter{
		Levels: []LogLevel{"DEBUG", "WARN", "ERROR"},
		ModifierFuncs: []ModifierFunc{
			Color(color.FgBlack),
			Color(color.FgYellow),
			Color(color.FgRed, color.Bold),
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
	expected := "\x1b[33m[WARN] foo\n\x1b[0m\x1b[31;1m[ERROR] bar\n\x1b[0m\x1b[33m[WARN] buzz\n\x1b[0m"
	if result != expected {
		t.Fatalf("bad: %#v", result)
	}
	t.Log(result)
}
