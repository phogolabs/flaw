package flaw

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// StackFrame represents a program counter inside a stack frame.
// For historical reasons if StackFrame is interpreted as a uintptr
// its value represents the program counter + 1.
type StackFrame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (frame StackFrame) pc() uintptr { return uintptr(frame) - 1 }

// file returns the full path to the file that contains the
// function for this StackFrame's pc.
func (frame StackFrame) file() string {
	fn := runtime.FuncForPC(frame.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(frame.pc())
	return file
}

// line returns the line number of source code of the
// function for this StackFrame's pc.
func (frame StackFrame) line() int {
	fn := runtime.FuncForPC(frame.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(frame.pc())
	return line
}

// name returns the name of this function, if known.
func (frame StackFrame) name() string {
	fn := runtime.FuncForPC(frame.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// Format formats the frame according to the fmt.Formatter interface.
//
//    %s    source file
//    %d    source line
//    %n    function name
//    %v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+s   function name and path of source file relative to the compile time
//          GOPATH separated by \n\t (<funcname>\n\t<path>)
//    %+v   equivalent to %+s:%d
func (frame StackFrame) Format(state fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case state.Flag('+'):
			fmt.Fprint(state, frame.file())
			fmt.Fprint(state, ":")
			frame.Format(state, 'd')
			fmt.Fprintf(state, " (%s)", frame.name())
		default:
			fmt.Fprint(state, path.Base(frame.file()))
			fmt.Fprint(state, ":")
			frame.Format(state, 'd')
		}
	case 'd':
		fmt.Fprint(state, strconv.Itoa(frame.line()))
	case 'n':
		fmt.Fprint(state, funcname(frame.name()))
	case 'v':
		switch {
		case state.Flag('+'):
			frame.Format(state, 's')
		default:
			fmt.Fprint(state, frame.file())
			fmt.Fprint(state, ":")
			frame.Format(state, 'd')
		}
	}
}

// MarshalText formats a stacktrace StackFrame as a text string. The output is the
// same as that of fmt.Sprintf("%+v", f), but without newlines or tabs.
func (frame StackFrame) MarshalText() ([]byte, error) {
	name := frame.name()

	if name == "unknown" {
		return []byte(name), nil
	}

	return []byte(fmt.Sprintf("%s %s:%d", name, frame.file(), frame.line())), nil
}

// StackTrace is stack of StackFrames from innermost (newest) to outermost (oldest).
type StackTrace []StackFrame

// NewStackTrace creates a new StackTrace
func NewStackTrace() StackTrace {
	const depth = 32

	var (
		pcs   [depth]uintptr
		trace stack
	)

	n := runtime.Callers(3, pcs[:])
	trace = pcs[0:n]

	return trace.StackTrace()
}

// Skip frames for given count
func (stack *StackTrace) Skip(n int) {
	count := len(*stack)

	if n > 0 && n < count {
		*stack = StackTrace((*stack)[n:])
	}
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
	case 'v':
		switch {
		case state.Flag('+'):
			for _, frame := range stack {
				fmt.Fprint(state, " --- ")
				frame.Format(state, verb)
				fmt.Fprint(state, "\n")
			}
		default:
			stack.formatSlice(state, verb)
		}
	case 's':
		stack.formatSlice(state, verb)
	}
}

// formatSlice will format this StackTrace into the given buffer as a slice of
// StackFrame, only valid when called with '%s' or '%v'.
func (stack StackTrace) formatSlice(state fmt.State, verb rune) {
	fmt.Fprint(state, "[")

	for index, frame := range stack {
		if index > 0 {
			fmt.Fprint(state, " ")
		}

		frame.Format(state, verb)
	}

	fmt.Fprint(state, "]")
}

// stack represents a stack of program counters.
type stack []uintptr

func (trace *stack) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case state.Flag('+'):
			for _, pc := range *trace {
				frame := StackFrame(pc)
				fmt.Fprintf(state, "\n%+v", frame)
			}
		}
	}
}

func (trace *stack) StackTrace() StackTrace {
	frames := make([]StackFrame, len(*trace))

	for index := 0; index < len(frames); index++ {
		frames[index] = StackFrame((*trace)[index])
	}

	return frames
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
