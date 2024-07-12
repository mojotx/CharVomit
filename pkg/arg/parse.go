package arg

import (
	"bytes"
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
	Version     bool
	Excluded    string
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

func Parse(fs *flag.FlagSet) (exitAfter bool, rc int) {

	var buffer bytes.Buffer
	flag.CommandLine.SetOutput(&buffer)
	output := flag.CommandLine.Output()

	fs.Usage = Usage
	fs.BoolVar(&Config.UpperCase, "u", false, "use upper-case letters")
	fs.BoolVar(&Config.LowerCase, "l", false, "use lower-case letters")
	fs.BoolVar(&Config.Digits, "d", false, "use numeric digits")
	fs.BoolVar(&Config.Symbols, "s", false, "use symbols: !#%+:=?@")
	fs.BoolVar(&Config.WeakChars, "w", false, "use weak characters (2-9, A-N, P-Z, a-k, m-z)")
	fs.BoolVar(&Config.ShowHelp, "h", false, "show help and exit")
	fs.BoolVar(&Config.Version, "v", false, "show version")
	fs.StringVar(&Config.Excluded, "x", "", "excluded characters (will be removed)")

	if err := fs.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(output, "cannot parse os.Args[1:]: %s\n", err.Error())
		exitAfter = true
		rc = 1
	} else if Config.Version {
		_, _ = fmt.Fprintln(output, Version())
		exitAfter = true
		rc = 0
	} else if fs.NArg() == 1 {
		Config.PasswordLen, err = strconv.Atoi(fs.Arg(0))
		if err != nil {
			_, _ = fmt.Fprintf(output, "cannot parse argument '%+v': %s\n", fs.Arg(0), err.Error())
			exitAfter = true
			rc = 1
		}

		// Get absolute value
		if Config.PasswordLen < 0 {
			Config.PasswordLen = Config.PasswordLen * -1
		}
	} else if Config.ShowHelp {
		fs.Usage()
		exitAfter = true
		rc = 0
	} else {
		// default to 32 characters
		Config.PasswordLen = 32
	}

	return
}
