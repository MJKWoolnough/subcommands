package subcommands

type options struct {
	summary, description string
	subcommand           []string
}

// Options represents an optional configuration setting to pass to the
// Registration Function.
type Option func(*options)

// Summary sets the subcommand specific summary shown in subcommand help text.
//
// To display a more comprehensive description in a subcommand specific help
// text, use the Description option.
func Summary(summary string) Option {
	return func(o *options) {
		o.summary = summary
	}
}

// Description sets the description presented in the subcommand specific
// help text.
func Description(desc string) Option {
	return func(o *options) {
		o.description = desc
	}
}

// Subcommand allows the given function to be registered to the given
// subcommand.
//
// Multiple parts can be specified, allow for sub-subcommands and similar.
func Subcommand(parts ...string) Option {
	return func(o *options) {
		o.subcommand = parts
	}
}
