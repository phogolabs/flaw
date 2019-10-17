package flaw

import (
	"fmt"
	"go/build"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// StackFrame represents a program counter inside a stack frame.
// For historical reasons if StackFrame is interpreted as a uintptr
// its value represents the program counter + 1.
type StackFrame runtime.Frame

// Format formats the frame according to the fmt.Formatter interface.
//
//    %s    source file
//    %d    source line
//    %n    function name
//    %v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+s   source file full path
//    %+v   equivalent to %+s:%d (%n)
func (frame StackFrame) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		frame.Format(state, 's')
		fmt.Fprintf(state, ":")
		frame.Format(state, 'd')

		if state.Flag('+') {
			fmt.Fprintf(state, " (")
			frame.Format(state, 'n')
			fmt.Fprintf(state, ")")
		}
	case 's':
		switch {
		case state.Flag('+'):
			fmt.Fprint(state, frame.File)
		default:
			path := frame.File

			if root := build.Default.GOPATH; root != "" {
				if file, err := filepath.Rel(root, path); err == nil {
					path = strings.TrimPrefix(file, "src/")
					path = strings.TrimPrefix(path, "pkg/mod/")
				}
			}

			fmt.Fprint(state, path)
		}
	case 'd':
		fmt.Fprint(state, strconv.Itoa(frame.Line))
	case 'n':
		name := frame.Function
		withoutPath := name[strings.LastIndex(name, "/")+1:]
		withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

		name = withoutPackage
		name = strings.Replace(name, "(", "", 1)
		name = strings.Replace(name, "*", "", 1)
		name = strings.Replace(name, ")", "", 1)

		fmt.Fprint(state, name)
	}
}

// MarshalText formats a stacktrace StackFrame as a text string. The output is the
// same as that of fmt.Sprintf("%+v", f), but without newlines or tabs.
func (frame StackFrame) MarshalText() ([]byte, error) {
	if name := frame.Function; name == "unknown" {
		return []byte(name), nil
	}

	return []byte(fmt.Sprintf("%+v", frame)), nil
}

// StackTrace is stack of StackFrames from innermost (newest) to outermost (oldest).
type StackTrace []StackFrame

// NewStackTrace creates a new StackTrace
func NewStackTrace() StackTrace {
	var (
		stack  = make([]uintptr, 32)
		size   = runtime.Callers(3, stack[:])
		frames = runtime.CallersFrames(stack[:size])
		trace  = StackTrace{}
	)

	for {
		frame, ok := frames.Next()
		if !ok {
			return trace
		}

		trace = append(trace, StackFrame(frame))
	}
}

// NewStackTraceAt creates a new stack trace at given position
func NewStackTraceAt(n int) StackTrace {
	n = n + 1
	stack := NewStackTrace()

	count := len(stack)

	if n > 0 && n < count {
		stack = StackTrace(stack[n:])
	}

	return stack
}

// Format formats the stack of StackFrames according to the fmt.Formatter interface.
//
//    %s	lists source files for each StackFrame in the stack
//    %v	lists the source file and line number for each StackFrame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+v   Prints filename, function, and line number for each StackFrame in the stack.
func (stack StackTrace) Format(state fmt.State, verb rune) {
	switch verb {
	case 's':
		fallthrough
	case 'v':
		switch {
		case state.Flag('+'):
			stack.formatBullet(state, verb)
		case state.Flag('#'):
			fmt.Fprintf(state, "%#v", []StackFrame(stack))
		default:
			stack.formatSlice(state, verb)
		}
	}
}

func (stack StackTrace) formatBullet(state fmt.State, verb rune) {
	count := len(stack)

	for index, frame := range stack {
		fmt.Fprint(state, " --- ")
		frame.Format(state, verb)

		if index < count-1 {
			fmt.Fprint(state, "\n")
		}
	}
}

func (stack StackTrace) formatSlice(state fmt.State, verb rune) {
	fmt.Fprint(state, "[")

	for index, frame := range stack {
		if index > 0 {
			fmt.Fprint(state, ", ")
		}

		frame.Format(state, verb)
	}

	fmt.Fprint(state, "]")
}
