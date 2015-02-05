package models

type arguments []string

func NewArguments(args []string) arguments {
	return arguments(args)
}
