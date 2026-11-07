package subcommands

type options struct {
	description, long string
	subcommand        []string
}

// Options represents an optional configuration setting to pass to the
// Registration Function.
type Option func(*options)

// Description sets the subcommand specific description shown in help text.
//
// To display a more comprehensive description in a subcommand specific help
// text, use the Long Description option.
func Description(desc string) Option {
	return func(o *options) {
		o.description = desc
	}
}

// LongDescription sets the description presented in the subcommand specific
// help text.
func LongDescription(desc string) Option {
	return func(o *options) {
		o.long = desc
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
