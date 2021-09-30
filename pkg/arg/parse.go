package arg

import (
	"flag"
	"fmt"
	"os"
)

type ConfigType struct {
	UpperCase bool
	LowerCase bool
	Digits    bool
	Symbols   bool
	WeakChars bool
	Args      int
}

var Config ConfigType

func Usage() {
	out := flag.CommandLine.Output()
	_, _ = fmt.Fprintf(out, "Usage: %s\n\n", os.Args[0])
	_, _ = fmt.Fprintf(out, "Options:\n")

	flag.PrintDefaults()

	_, _ = fmt.Fprintln(out)

}

func Parse() {

	flag.Usage = Usage
	flag.BoolVar(&Config.UpperCase, "uc", false, "use upper-case letters")
	flag.BoolVar(&Config.LowerCase, "lc", false, "use lower-case letters")
	flag.BoolVar(&Config.UpperCase, "d", false, "use numeric digits")
	flag.BoolVar(&Config.UpperCase, "s", false, "use symbols")
	flag.BoolVar(&Config.UpperCase, "w", false, "use weak characters (2-9, A-Z, a-z)")

	Config.Args = flag.NArg()

}
