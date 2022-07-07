package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

const workdingDir = "/data/local/tmp/odexpatcher"

type CommandOptions struct {
	VersionFlag bool     `long:"version" description:"Show the app version information"`
	Patch       PatchCmd `command:"patch" description:"Generate patched oat from dex"`
	Info        InfoCmd  `command:"info" description:"Get oat file info"`
}

func main() {
	opts := &CommandOptions{}
	parser := flags.NewParser(opts, flags.Default)
	parser.SubcommandsOptional = true
	_, err := parser.Parse()
	if err != nil {
		switch flagsErr := err.(type) {
		case *flags.Error:
			if flagsErr.Type == flags.ErrHelp {
				os.Exit(0)
			}
		default:
			log.Fatalln(err)
		}
	}

	if opts.VersionFlag {
		printVersion()
		os.Exit(0)
	}

}
