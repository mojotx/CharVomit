package arg

import (
	"flag"
	"fmt"
	"os"
)

type ConfigType struct {
	Args      int
	Digits    bool
	ShowHelp  bool
	LowerCase bool
	Symbols   bool
	UpperCase bool
	WeakChars bool
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
	flag.BoolVar(&Config.UpperCase, "u", false, "use upper-case letters")
	flag.BoolVar(&Config.LowerCase, "l", false, "use lower-case letters")
	flag.BoolVar(&Config.Digits, "d", false, "use numeric digits")
	flag.BoolVar(&Config.Symbols, "s", false, "use symbols: !#%+:=?@")
	flag.BoolVar(&Config.WeakChars, "w", false, "use weak characters (2-9, A-Z, a-z)")
	flag.BoolVar(&Config.ShowHelp, "h", false, "show help and exit")

	flag.Parse()

	Config.Args = flag.NArg()

}
