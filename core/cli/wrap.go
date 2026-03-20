package cli

import "flag"

type FlagFn func(fs *flag.FlagSet, args []string) error

func FlagWrap(args []string, flagFn FlagFn) error {
	fs := flag.NewFlagSet(args[1], flag.ContinueOnError)
	return flagFn(fs, args[2:])
}

type DescFn func() string