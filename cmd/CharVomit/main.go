package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

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

	pwLen := 32

	if fs.NArg() == 1 {

		var err error
		pwLen, err = strconv.Atoi(fs.Arg(0))
		if err != nil {
			fmt.Printf("cannot parse argument '%+v': %s", flag.Arg(0), err.Error())
			os.Exit(1)
		}

		// Get absolute value
		if pwLen < 0 {
			pwLen = pwLen * -1
		}
	}

	pw, err := cv.Puke(pwLen)
	if err != nil {
		fmt.Printf("Puke(%d) error: %s", pwLen, err.Error())
		os.Exit(1)
	}

	fmt.Println(pw)
}
