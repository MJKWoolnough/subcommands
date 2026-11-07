package subcommands

type options struct {
	description, long string
	subcommand        []string
}

type Option func(*options)

func Description(desc string) Option {
	return func(o *options) {
		o.description = desc
	}
}

func LongDescription(desc string) Option {
	return func(o *options) {
		o.long = desc
	}
}

func Subcommand(parts ...string) Option {
	return func(o *options) {
		o.subcommand = parts
	}
}
