package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/minor-industries/packager/pkg/packager"
	"github.com/pkg/errors"
	"os"
)

func run() error {
	opts := new(packager.Opts)

	args, err := flags.Parse(opts)
	if err != nil {
		return err
	}

	name := args[0]

	if !opts.Minor {
		return errors.New("only minor releases supported for now")
	}

	switch opts.Arch {
	case "arm64":
	default:
		return errors.New("unknown arch")
	}

	opts.SharedFolder = os.ExpandEnv(opts.SharedFolder)

	return packager.Run(name, opts, buildSingle)
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
