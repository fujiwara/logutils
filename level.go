// Package logutils augments the standard log package with levels.
package logutils

import (
	"bytes"
	"io"
	"sync"
)

type ModifierFunc func([]byte) []byte

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

	// ModifierFuncs is the list of modifier functions to apply to each log lines.
	ModifierFuncs []ModifierFunc

	// MinLevel is the minimum level allowed through
	MinLevel LogLevel

	// The underlying io.Writer where log messages that pass the filter
	// will be set.
	Writer io.Writer

	modifierFuncs map[LogLevel]ModifierFunc
	once          sync.Once
}

// Check will check a given line if it would be included in the level
// filter.
func (f *LevelFilter) Check(line []byte) bool {
	_, ok := f.getModifierFunc(line)
	return ok
}

func (f *LevelFilter) getModifierFunc(line []byte) (ModifierFunc, bool) {
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

	mf, ok := f.modifierFuncs[level]
	return mf, ok
}

func (f *LevelFilter) Write(p []byte) (n int, err error) {
	// Note in general that io.Writer can receive any byte sequence
	// to write, but the "log" package always guarantees that we only
	// get a single line. We use that as a slight optimization within
	// this method, assuming we're dealing with a single, complete line
	// of log data.

	mf, ok := f.getModifierFunc(p)
	if !ok {
		// disabled log level
		return len(p), nil
	}
	if mf != nil {
		return f.Writer.Write(mf(p))
	} else {
		// default
		return f.Writer.Write(p)
	}
}

// SetMinLevel is used to update the minimum log level
func (f *LevelFilter) SetMinLevel(min LogLevel) {
	f.MinLevel = min
	f.init()
}

func (f *LevelFilter) init() {
	mfuncs := make(map[LogLevel]ModifierFunc, len(f.Levels))
	minLevelIndex := -1
	for i, level := range f.Levels {
		if i < len(f.ModifierFuncs) {
			mfuncs[level] = f.ModifierFuncs[i]
		} else {
			mfuncs[level] = nil
		}
		if level == f.MinLevel {
			minLevelIndex = i
		}
	}
	f.modifierFuncs = mfuncs
	if minLevelIndex == -1 {
		return
	}
	for i, level := range f.Levels {
		if i < minLevelIndex {
			delete(f.modifierFuncs, level)
		} else {
			return
		}
	}
}
