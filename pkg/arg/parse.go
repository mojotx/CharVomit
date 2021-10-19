package arg

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type ConfigType struct {
	PasswordLen int
	Digits      bool
	ShowHelp    bool
	LowerCase   bool
	Symbols     bool
	UpperCase   bool
	WeakChars   bool
}

var Config ConfigType

func Usage() {
	out := flag.CommandLine.Output()
	_, _ = fmt.Fprintf(out, "Usage: %s [ length ]\n\n", os.Args[0])
	_, _ = fmt.Fprintf(out, "If a password length is not specified, 32 is used.\n\n")
	_, _ = fmt.Fprintf(out, "Other optional flags are:\n")

	flag.PrintDefaults()

	_, _ = fmt.Fprintln(out)
	_, _ = fmt.Fprintf(out, "Note that optional flags must precede the password length.\n\n")
	_, _ = fmt.Fprintf(out, "For example, a 8-character password of all capital letters:\n")
	_, _ = fmt.Fprintf(out, "%s -u 8\n\n", os.Args[0])
	_, _ = fmt.Fprintln(out, "Also note that certain characters that are confusing are ignored by default,")
	_, _ = fmt.Fprintln(out, "such as '0', 'O', '1', and 'l'. You can still get those characters, if you wish,")
	_, _ = fmt.Fprintln(out, "by using the -u, -l, and -d flags. The default is equivalent to -w -s.")

}

func Parse() {
	flag.Usage = Usage
	flag.BoolVar(&Config.UpperCase, "u", false, "use upper-case letters")
	flag.BoolVar(&Config.LowerCase, "l", false, "use lower-case letters")
	flag.BoolVar(&Config.Digits, "d", false, "use numeric digits")
	flag.BoolVar(&Config.Symbols, "s", false, "use symbols: !#%+:=?@")
	flag.BoolVar(&Config.WeakChars, "w", false, "use weak characters (2-9, A-N, P-Z, a-k, m-z)")
	flag.BoolVar(&Config.ShowHelp, "h", false, "show help and exit")

	flag.Parse()

	if flag.NArg() == 1 {

		var err error
		Config.PasswordLen, err = strconv.Atoi(flag.Arg(0))
		if err != nil {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), "cannot parse argument '%+v': %s", flag.Arg(0), err.Error())
			os.Exit(1)
		}

		// Get absolute value
		if Config.PasswordLen < 0 {
			Config.PasswordLen = Config.PasswordLen * -1
		}
	} else {
		// default to 32 characters
		Config.PasswordLen = 32
	}

}
