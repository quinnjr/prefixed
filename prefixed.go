// Package prefixed extends https://github.com/apex/log/handlers/cli with a prefix.
package prefixed

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

// Default handler outputting to stderr.
var Default = New("", os.Stderr)

// Handler implementation.
type Handler struct {
	mu      sync.Mutex
	Writer  io.Writer
	Padding int
	Prefix  string
}

// New handler.
func New(p string, w io.Writer) *Handler {
	return &Handler{
		Writer:  w,
		Padding: 3,
		Prefix:  p,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	color := cli.Colors[e.Level]
	level := cli.Strings[e.Level]
	names := e.Fields.Names()
	e.Message = fmt.Sprintf("%s = %s", h.Prefix, e.Message)

	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Fprintf(h.Writer, "\033[%dm%*s\033[0m %-25s", color, h.Padding+1, level, e.Message)

	for _, name := range names {
		if name == "source" {
			continue
		}

		fmt.Fprintf(h.Writer, " \033[%dm%s\033[0m=%v", color, name, e.Fields.Get(name))
	}

	fmt.Fprintln(h.Writer)

	return nil
}
