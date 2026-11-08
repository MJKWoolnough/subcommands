package subcommands

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"unicode/utf8"
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

const spacing = 2

func (h *Help) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	if h.Err != nil {
		fmt.Fprintf(&buf, "Error: %v\n\n", h.Err)
	}

	exe, err := os.Executable()
	if err != nil {
		return 0, err
	}

	fmt.Fprintf(&buf, "Usage: %s", exe)

	if len(h.Subcommand) > 0 {
		for _, sub := range h.Subcommand {
			fmt.Fprintf(&buf, " %s", sub)
		}
	}

	flags := h.flags.list()

	for _, flag := range flags {
		fmt.Fprintf(&buf, " %s", flag.aliases[0])

		if flag.typ.Implements(flagI) || reflect.PointerTo(flag.typ).Implements(flagI) {
			fmt.Fprint(&buf, " value")
		} else if flag.typ.Kind() != reflect.Bool {
			fmt.Fprintf(&buf, " %s", flag.typ.Kind())
		}
	}

	if desc := cmp.Or(h.command.long, h.command.short); desc != "" {
		fmt.Fprint(&buf, "\n\n", desc)
	}

	if len(flags) > 0 {
		fmt.Fprint(&buf, "\n\nFlags\n=====\n")

		maxLen := 0

		for _, flag := range flags {
			maxLen = max(maxLen, flag.aliasLength())
		}

		for _, flag := range flags {
			io.WriteString(&buf, "\n")

			for n, alias := range flag.aliases {
				if n > 0 {
					io.WriteString(&buf, ",")
				}

				io.WriteString(&buf, alias)
			}

			io.WriteString(&buf, strings.Repeat(" ", maxLen-flag.aliasLength()+spacing))
			io.WriteString(&buf, flag.usage)
		}
	}

	if len(h.subcommands) > 0 {
		fmt.Fprint(&buf, "\n\nSubcommands\n===========\n")

		cmds := h.list()

		maxLen := 0

		for _, cmd := range cmds {
			maxLen = max(maxLen, utf8.RuneCountInString(cmd.name))
		}

		for _, cmd := range cmds {
			io.WriteString(&buf, "\n")
			io.WriteString(&buf, cmd.name)
			io.WriteString(&buf, strings.Repeat(" ", maxLen-len(cmd.name)+spacing))
			io.WriteString(&buf, cmd.short)
		}
	}

	return buf.WriteTo(w)
}
