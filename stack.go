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

func (frame StackFrame) short() string {
	name := frame.name()
	withoutPath := name[strings.LastIndex(name, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
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
			fmt.Fprint(state, frame.file())
		default:
			fmt.Fprint(state, path.Base(frame.file()))
		}
	case 'd':
		fmt.Fprint(state, strconv.Itoa(frame.line()))
	case 'n':
		fmt.Fprint(state, frame.short())
	}
}

// MarshalText formats a stacktrace StackFrame as a text string. The output is the
// same as that of fmt.Sprintf("%+v", f), but without newlines or tabs.
func (frame StackFrame) MarshalText() ([]byte, error) {
	name := frame.name()

	if name == "unknown" {
		return []byte(name), nil
	}

	text := fmt.Sprintf("%s:%d (%s)", frame.file(), frame.line(), name)
	return []byte(text), nil
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
