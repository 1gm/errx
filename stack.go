package errx

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

var callerSkipLevel int

// SetCallerSkipLevel sets the number of callers to skip when building a stack frame.
// By default it is set to 0, causing all stack frames to originate at the point where
// errx.New, errx.Errorf, errx.Wrap, or errx.Wrapf was called.
//
// If any of these functions is wrapped, SetCallerSkipLevel should be called with the
// number of wrapping functions.
//
// For example, if errx.New was wrapped in a helper function, e.g. SetupError(args...),
// then SetCallerSkipLevel should be called with a value of 1.
//
// SetCallerSkipLevel should only be called once. It is not goroutine safe and is intended
// to be set as part of an initialization routine.
func SetCallerSkipLevel(level int) {
	callerSkipLevel = level
}

type StackFrame struct {
	Name            string
	File            string
	TrimmedFileLine string
	Line            int
}

func (s *StackFrame) String() string {
	return fmt.Sprintf("  at %s(%s:%d)\n", s.Name, s.TrimmedFileLine, s.Line)
}

type StackTrace struct {
	Frames []StackFrame
}

func (s *StackTrace) String() string {
	b := new(bytes.Buffer)
	for _, f := range s.Frames {
		b.WriteString(f.String())
	}
	return b.String()
}

const maxStackDepth = 32

func getStack() *StackTrace {
	st := &StackTrace{}

	var pcs [maxStackDepth]uintptr
	n := runtime.Callers(4+callerSkipLevel, pcs[:])
	for _, pc := range pcs[0:n] {
		pcFunc := runtime.FuncForPC(pc)
		name := pcFunc.Name()
		file, line := pcFunc.FileLine(pc)
		trimmed := trimGOPATH(name, file)
		st.Frames = append(st.Frames, StackFrame{name, file, trimmed, line})
	}
	return st
}

const (
	sep    = "/"
	sepLen = len(sep)
)

// This code is taken from pkg/errors and modified very lightly.
// Copyright (c) 2015, Dave Cheney <dave@cheney.net>
func trimGOPATH(name, file string) string {
	// Here we want to get the source file path relative to the compile time
	// GOPATH. As of Go 1.6.x there is no direct way to know the compiled
	// GOPATH at runtime, but we can infer the number of path segments in the
	// GOPATH. We note that fn.Name() returns the function name qualified by
	// the import path, which does not include the GOPATH. Thus we can trim
	// segments from the beginning of the file path until the number of path
	// separators remaining is one more than the number of path separators in
	// the function name. For example, given:
	//
	//    GOPATH     /home/user
	//    file       /home/user/src/pkg/sub/file.go
	//    fn.Name()  pkg/sub.Type.Method
	//
	// We want to produce:
	//
	//    pkg/sub/file.go
	//
	// From this we can easily see that fn.Name() has one less path separator
	// than our desired output. We count separators from the end of the file
	// path until it finds two more than in the function name and then move
	// one character forward to preserve the initial path segment without a
	// leading separator.
	goal := strings.Count(name, sep) + 2
	i := len(file)
	for n := 0; n < goal; n++ {
		i = strings.LastIndex(file[:i], sep)
		if i == -1 {
			// not enough separators found, set i so that the slice expression
			// below leaves file unmodified
			i = -sepLen
			break
		}
	}
	// get back to 0 or trim the leading separator
	file = file[i+sepLen:]
	return file
}
