package format

import (
	"fmt"
	"io"
	"text/tabwriter"
)

var _ fmt.State = &State{}

// State wraps the state into desired writer
type State struct {
	state  fmt.State
	writer io.Writer
	size   int
}

// NewState creates a new state
func NewState(state fmt.State) *State {
	wstate := &State{
		state:  state,
		writer: state,
	}

	if state.Flag('+') {
		wstate.writer = tabwriter.NewWriter(state, 0, 0, 1, ' ', tabwriter.AlignRight)
	}

	return wstate
}

// Write is the function to call to emit formatted output to be printed.
func (w *State) Write(data []byte) (n int, err error) {
	n, err = w.writer.Write(data)

	if err == nil {
		w.size = w.size + n
	}

	return n, err
}

// Width returns the value of the width option and whether it has been set.
func (w *State) Width() (wid int, ok bool) {
	return w.state.Width()
}

// Precision returns the value of the precision option and whether it has been set.
func (w *State) Precision() (prec int, ok bool) {
	return w.state.Precision()
}

// Flag reports whether the flag c, a character, has been set.
func (w *State) Flag(c int) bool {
	return w.state.Flag(c)
}

// Flush flushes the state
func (w *State) Flush() error {
	type Flusher interface {
		Flush() error
	}

	if flusher, ok := w.writer.(Flusher); ok {
		return flusher.Flush()
	}

	return nil
}

// Size returns the size
func (w *State) Size() int {
	return w.size
}

// StringSlice represents a slice of string
type StringSlice []string

// Format formats the string slice
func (d StringSlice) Format(state fmt.State, verb rune) {
	switch verb {
	case 's':
		fallthrough
	case 'v':
		switch {
		case state.Flag('+'):
			d.formatBullet(state, verb)
		case state.Flag('#'):
			fmt.Fprintf(state, "%#v", []string(d))
		default:
			d.formatSlice(state, verb)
		}
	}
}

func (d StringSlice) formatBullet(state fmt.State, verb rune) {
	count := len(d)

	for index, line := range d {
		fmt.Fprint(state, " --- ")
		fmt.Fprint(state, line)

		if index < count-1 {
			fmt.Fprint(state, "\n")
		}
	}
}

func (d StringSlice) formatSlice(state fmt.State, verb rune) {
	fmt.Fprint(state, "[")

	for index, line := range d {
		if index > 0 {
			fmt.Fprint(state, ", ")
		}

		fmt.Fprint(state, line)
	}

	fmt.Fprint(state, "]")
}
