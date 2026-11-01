package subcommands

type (
	Required  struct{}
	Mutual    struct{}
	Exclusive struct{}
	Requires  struct{}
	Forbids   struct{}
	OneOf     struct{}
	AnyOf     struct{}
)

type Commands struct {}

var globalCommand Commands

func (c *Commands) Register[T any](short, long string, fn func(T) error, subs ...string) error {
	return nil
}

func (c *Commands) Prepare[T, U any](fn func(T) (U, error)) error {
	return nil
}

func (c *Commands) Run() error {
	return nil
}

func Register[T any](short, long string, fn func(T) error, subs ...string) error {
	return globalCommand.Register(short, long, fn, subs...)
}

func Prepare[T, U any](fn func(T) (U, error)) error {
	return globalCommand.Prepare(fn)
}

func Run() error {
	return c.Run()
}
