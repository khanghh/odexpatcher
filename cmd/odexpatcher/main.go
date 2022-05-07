package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

type CommandOptions struct {
	VersionFlag  bool    `long:"version" description:"Show the app version information"`
	Exchange     string  `long:"exchange" description:"Crypto currencies exchange"`
	TargetToken  string  `long:"target" description:"Snipping tareget token"`
	SourceToken  string  `long:"source" description:"Source token to swap from (leave empty to use native currency)"`
	BuyAmountIn  float64 `long:"buyAmountIn" descript:"Amount of source token"`
	BuyAmountOut float64 `long:"buyAmountOut" descript:"Amount of target token"`
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

	fmt.Println("aaa")

}
