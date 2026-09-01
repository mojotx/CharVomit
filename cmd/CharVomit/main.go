package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mojotx/CharVomit/pkg/CharVomit"
	"github.com/mojotx/CharVomit/pkg/arg"
)

var fs *flag.FlagSet

func init() {
	fs = flag.NewFlagSet("DynamicParser", flag.ContinueOnError)
}

// TO-DO:
// - Add support for duplicate character checking
func main() {
	cfg, shouldExit, rc := arg.ParseArgs(os.Args[1:], fs)
	if shouldExit {
		os.Exit(rc)
	}

	var cv CharVomit.CharVomit

	if err := cv.SetAcceptableChars(cfg); err != nil {
		fmt.Printf("could not set acceptable chars: %s\n", err.Error())
		os.Exit(1)
	}

	pw, err := cv.Puke(cfg.PasswordLen)
	if err != nil {
		fmt.Printf("Puke(%d) error: %s", cfg.PasswordLen, err.Error())
		os.Exit(1)
	}

	fmt.Println(pw)
}
