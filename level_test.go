package logutils

import (
	"bytes"
	"io"
	"log"
	"testing"
)

func TestLevelFilter_impl(t *testing.T) {
	var _ io.Writer = new(LevelFilter)
}

func TestLevelFilter(t *testing.T) {
	buf := new(bytes.Buffer)
	filter := &LevelFilter{
		Levels:   []LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "WARN",
		Writer:   buf,
	}

	logger := log.New(filter, "", 0)
	logger.Print("[WARN] foo")
	logger.Println("[ERROR] bar")
	logger.Println("[DEBUG] baz")
	logger.Println("[WARN] buzz")
	logger.Println("foobarbaz")
	logger.Println("[xxxx] foobarbaz")
	logger.Println(`{"foo":["bar","baz"]}`)

	result := buf.String()
	expected := "[WARN] foo\n[ERROR] bar\n[WARN] buzz\nfoobarbaz\n[xxxx] foobarbaz\n{\"foo\":[\"bar\",\"baz\"]}\n"
	if result != expected {
		t.Fatalf("expected: %#v, bad: %#v", expected, result)
	}
}

func TestLevelFilterWithPrefix(t *testing.T) {
	buf := new(bytes.Buffer)
	filter := &LevelFilter{
		Levels:   []LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "WARN",
		Writer:   buf,
	}

	logger := log.New(filter, "prefix ", 0)
	logger.Print("[WARN] foo")
	logger.Println("[ERROR] bar")
	logger.Println("[DEBUG] baz")
	logger.Println("[WARN] buzz")
	logger.Println("foobarbaz")
	logger.Println("[xxxx] foobarbaz")
	logger.Println(`{"foo":["bar","baz"]}`)

	result := buf.String()
	expected := "prefix [WARN] foo\nprefix [ERROR] bar\nprefix [WARN] buzz\nprefix foobarbaz\nprefix [xxxx] foobarbaz\nprefix {\"foo\":[\"bar\",\"baz\"]}\n"
	if result != expected {
		t.Fatalf("expected: %#v, bad: %#v", expected, result)
	}
}

func TestLevelFilterCheck(t *testing.T) {
	filter := &LevelFilter{
		Levels:   []LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "WARN",
		Writer:   nil,
	}

	testCases := []struct {
		line  string
		check bool
	}{
		{"[WARN] foo\n", true},
		{"[ERROR] bar\n", true},
		{"[DEBUG] baz\n", false},
		{"[WARN] buzz\n", true},
	}

	for _, testCase := range testCases {
		result := filter.Check([]byte(testCase.line))
		if result != testCase.check {
			t.Errorf("Fail: %s", testCase.line)
		}
	}
}

func TestLevelFilter_SetMinLevel(t *testing.T) {
	filter := &LevelFilter{
		Levels:   []LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "ERROR",
		Writer:   nil,
	}

	testCases := []struct {
		line        string
		checkBefore bool
		checkAfter  bool
	}{
		{"[WARN] foo\n", false, true},
		{"[ERROR] bar\n", true, true},
		{"[DEBUG] baz\n", false, false},
		{"[WARN] buzz\n", false, true},
	}

	for _, testCase := range testCases {
		result := filter.Check([]byte(testCase.line))
		if result != testCase.checkBefore {
			t.Errorf("Fail: %s", testCase.line)
		}
	}

	// Update the minimum level to WARN
	filter.SetMinLevel("WARN")

	for _, testCase := range testCases {
		result := filter.Check([]byte(testCase.line))
		if result != testCase.checkAfter {
			t.Errorf("Fail: %s", testCase.line)
		}
	}
}
