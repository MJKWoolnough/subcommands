package subcommands

type prep struct {
	flags, result reflect.Type
	fn            reflect.Value
}
