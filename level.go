// Package logutils augments the standard log package with levels.
package logutils

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/fatih/color"
)

type LogLevel string

// LevelFilter is an io.Writer that can be used with a logger that
// will filter out log messages that aren't at least a certain level.
//
// Once the filter is in use somewhere, it is not safe to modify
// the structure.
type LevelFilter struct {
	// Levels is the list of log levels, in increasing order of
	// severity. Example might be: {"DEBUG", "WARN", "ERROR"}.
	Levels []LogLevel

	// Colors is the list of github.com/fatih/color.Color to use for each log level.
	Colors []color.Color

	// MinLevel is the minimum level allowed through
	MinLevel LogLevel

	// The underlying io.Writer where log messages that pass the filter
	// will be set.
	Writer io.Writer

	printers map[LogLevel]printer
	once     sync.Once
}

type printer func(io.Writer, ...interface{}) (int, error)

// Check will check a given line if it would be included in the level
// filter.
func (f *LevelFilter) Check(line []byte) bool {
	return f.Printer(line) != nil
}

func (f *LevelFilter) Printer(line []byte) printer {
	f.once.Do(f.init)

	// Check for a log level
	var level LogLevel
	x := bytes.IndexByte(line, '[')
	if x >= 0 {
		y := bytes.IndexByte(line[x:], ']')
		if y >= 0 {
			level = LogLevel(line[x+1 : x+y])
		}
	}

	return f.printers[level]
}

func (f *LevelFilter) Write(p []byte) (n int, err error) {
	// Note in general that io.Writer can receive any byte sequence
	// to write, but the "log" package always guarantees that we only
	// get a single line. We use that as a slight optimization within
	// this method, assuming we're dealing with a single, complete line
	// of log data.

	if pr := f.Printer(p); pr != nil {
		return pr(f.Writer, string(p))
	}
	return len(p), nil
}

// SetMinLevel is used to update the minimum log level
func (f *LevelFilter) SetMinLevel(min LogLevel) {
	f.MinLevel = min
	f.init()
}

func (f *LevelFilter) init() {
	printers := make(map[LogLevel]printer, len(f.Levels))
	minLevelIndex := -1
	for i, level := range f.Levels {
		if i < len(f.Colors) {
			printers[level] = f.Colors[i].Fprint
		} else {
			printers[level] = fmt.Fprint
		}
		if level == f.MinLevel {
			minLevelIndex = i
		}
	}
	f.printers = printers
	if minLevelIndex == -1 {
		return
	}
	for i, level := range f.Levels {
		if i < minLevelIndex {
			delete(f.printers, level)
		} else {
			return
		}
	}
}
