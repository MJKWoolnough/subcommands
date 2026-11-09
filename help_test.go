package subcommands

import (
	"io"
	"os"
	"testing"
)

func TestHelpWriteTo(t *testing.T) {
	exe, _ := os.Executable()

	for n, test := range [...]struct {
		Commands   func(*Commands)
		Subcommand []string
		Err        error
		Output     string
	}{
		{
			func(*Commands) {},
			nil,
			nil,
			"Usage: " + exe,
		},
		{
			func(*Commands) {},
			nil,
			io.EOF,
			"Error: EOF\n\nUsage: " + exe,
		},
		{
			func(c *Commands) {
				c.Register(func(struct{}) error { return nil }, Description("My Description"))
			},
			nil,
			nil,
			"Usage: " + exe + "\n\nMy Description",
		},
		{
			func(c *Commands) {
				c.Register(func(struct {
					A  int
					Bd string
				},
				) error {
					return nil
				})
			},
			nil,
			nil,
			"Usage: " + exe + " --a int --bd string\n\nFlags\n=====\n\n--a\n--bd",
		},
	} {
		var c Commands

		test.Commands(&c)

		sc := &c.command

		for _, sub := range test.Subcommand {
			sc = sc.subcommands[sub]
		}

		h := Help{command: sc, Subcommand: test.Subcommand, Err: test.Err}

		if ht := h.Error(); ht != test.Output {
			t.Errorf("test %d: expecting output %q, got %q", n+1, test.Output, ht)
		}
	}
}
