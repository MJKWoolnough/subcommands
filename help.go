package subcommands

import (
	"io"
	"strings"
)

type Help struct {
	*command
	Subcommand []string
	Err        error
}

func (h *Help) Error() string {
	var buf strings.Builder

	h.WriteTo(&buf)

	return buf.String()
}

func (h *Help) Unwrap() error {
	return h.Err
}

func (h *Help) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}
