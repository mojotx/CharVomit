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
	fs = flag.NewFlagSet("DynamicParser", flag.ExitOnError)
}

// TO-DO:
// - Add support for duplicate character checking
func main() {

	shouldExit, rc := arg.Parse(fs)
	if shouldExit {
		os.Exit(rc)
	}

	if arg.Config.ShowHelp {
		arg.Usage()
		os.Exit(0)
	}

	var cv CharVomit.CharVomit

	if err := cv.SetAcceptableChars(arg.Config); err != nil {
		fmt.Printf("could not set acceptable chars: %s\n", err.Error())
		os.Exit(1)
	}

	pw, err := cv.Puke(arg.Config.PasswordLen)
	if err != nil {
		fmt.Printf("Puke(%d) error: %s", arg.Config.PasswordLen, err.Error())
		os.Exit(1)
	}

	fmt.Println(pw)
}
