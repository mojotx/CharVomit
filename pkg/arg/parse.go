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
	Stdout      bool
	OutputFile  string
	Version     bool
	Excluded    string
}

var Config ConfigType

func Usage() {
	out := flag.CommandLine.Output()
	_, _ = fmt.Fprintf(out, "Usage: %s [ length ]\n\n", os.Args[0])
	_, _ = fmt.Fprintf(out, "If a password length is not specified, 32 is used.\n")
	_, _ = fmt.Fprintln(out, "With no character flags, the default pool is equivalent to -w -s: weak, non-ambiguous characters plus symbols.")
	_, _ = fmt.Fprintln(out, "The -d flag is redundant with -w because the weak pool already contains the non-ambiguous digits.")
	_, _ = fmt.Fprintln(out)
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

// makeUsageFunc creates a usage function for the given FlagSet
func makeUsageFunc(fs *flag.FlagSet) func() {
	return func() {
		out := fs.Output()
		_, _ = fmt.Fprintf(out, "Usage: %s [ length ]\n\n", os.Args[0])
		_, _ = fmt.Fprintf(out, "If a password length is not specified, 32 is used.\n")
		_, _ = fmt.Fprintln(out, "With no character flags, the default pool is equivalent to -w -s: weak, non-ambiguous characters plus symbols.")
		_, _ = fmt.Fprintln(out, "The -d flag is redundant with -w because the weak pool already contains the non-ambiguous digits.")
		_, _ = fmt.Fprintln(out)
		_, _ = fmt.Fprintf(out, "Other optional flags are:\n")

		fs.PrintDefaults()

		_, _ = fmt.Fprintln(out)
		_, _ = fmt.Fprintf(out, "Note that optional flags must precede the password length.\n\n")
		_, _ = fmt.Fprintf(out, "For example, a 8-character password of all capital letters:\n")
		_, _ = fmt.Fprintf(out, "%s -u 8\n\n", os.Args[0])
		_, _ = fmt.Fprintln(out, "Also note that certain characters that are confusing are ignored by default,")
		_, _ = fmt.Fprintln(out, "such as '0', 'O', '1', and 'l'. You can still get those characters, if you wish,")
		_, _ = fmt.Fprintln(out, "by using the -u, -l, and -d flags. The default is equivalent to -w -s.")
	}
}

// RegisterFlags registers all command-line flags with the given FlagSet and target config.
func RegisterFlags(fs *flag.FlagSet, cfg *ConfigType) {
	fs.BoolVar(&cfg.UpperCase, "u", false, "use upper-case letters")
	fs.BoolVar(&cfg.LowerCase, "l", false, "use lower-case letters")
	fs.BoolVar(&cfg.Digits, "d", false, "use numeric digits")
	fs.BoolVar(&cfg.Symbols, "s", false, "use symbols: !#%+:=?@")
	fs.BoolVar(&cfg.WeakChars, "w", false, "use weak characters (2-9, A-N, P-Z, a-k, m-z)")
	fs.BoolVar(&cfg.ShowHelp, "h", false, "show help and exit")
	fs.BoolVar(&cfg.Version, "v", false, "show version")
	fs.StringVar(&cfg.Excluded, "x", "", "excluded characters (will be removed from char pool)")
}

func resetFlagSetDefaults(fs *flag.FlagSet) {
	fs.VisitAll(func(f *flag.Flag) {
		_ = fs.Set(f.Name, f.DefValue)
	})
}

// ParseArgs parses CLI arguments into a ConfigType without relying on package-level state.
func ParseArgs(args []string, fs *flag.FlagSet) (cfg ConfigType, exitAfter bool, rc int) {
	Config = ConfigType{}
	output := fs.Output()
	fs.Usage = makeUsageFunc(fs)
	if fs.Lookup("u") == nil {
		RegisterFlags(fs, &Config)
	}
	resetFlagSetDefaults(fs)

	if err := fs.Parse(args); err != nil {
		_, _ = fmt.Fprintf(output, "cannot parse os.Args[1:]: %s\n", err.Error())
		cfg = Config
		return cfg, true, 1
	}

	cfg = Config
	if cfg.Version {
		_, _ = fmt.Fprintln(output, Version())
		return cfg, true, 0
	}
	if cfg.ShowHelp {
		fs.Usage()
		return cfg, true, 0
	}
	if fs.NArg() > 1 {
		_, _ = fmt.Fprintf(output, "too many arguments: expected at most 1 positional length, got %d\n", fs.NArg())
		return cfg, true, 1
	}

	cfg.PasswordLen = 32
	if fs.NArg() == 1 {
		parsedLen, err := strconv.Atoi(fs.Arg(0))
		if err != nil {
			_, _ = fmt.Fprintf(output, "cannot parse argument '%+v': %s\n", fs.Arg(0), err.Error())
			return cfg, true, 1
		}

		cfg.PasswordLen = parsedLen
		if cfg.PasswordLen < 0 {
			cfg.PasswordLen *= -1
		}
	}

	Config = cfg
	return cfg, false, 0
}

// ParseConfig is kept as a compatibility wrapper around ParseArgs.
func ParseConfig(fs *flag.FlagSet) (cfg ConfigType, exitAfter bool, rc int) {
	cfg, exitAfter, rc = ParseArgs(os.Args[1:], fs)
	Config = cfg
	return cfg, exitAfter, rc
}

func Parse(fs *flag.FlagSet) (exitAfter bool, rc int) {
	cfg, exitAfter, rc := ParseArgs(os.Args[1:], fs)
	Config = cfg
	return exitAfter, rc
}
