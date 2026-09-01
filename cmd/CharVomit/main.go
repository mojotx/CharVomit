package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/mojotx/CharVomit/pkg/CharVomit"
	"github.com/mojotx/CharVomit/pkg/arg"
)

var rootCmd = &cobra.Command{
	Use:           "CharVomit [length]",
	Short:         "Generate random passwords.",
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := resolveCommandConfig(cmd, args)
		if err != nil {
			return err
		}

		return run(cfg)
	},
}

func init() {
	rootCmd.Flags().StringP("output-file", "o", "", "write the generated password to a file instead of stdout")
	rootCmd.Flags().BoolP("digits", "d", false, "use numeric digits")
	rootCmd.Flags().BoolP("lowercase", "l", false, "use lower-case letters")
	rootCmd.Flags().BoolP("uppercase", "u", false, "use upper-case letters")
	rootCmd.Flags().BoolP("symbols", "s", true, "use symbols: !#%+:=?@")
	rootCmd.Flags().BoolP("weak", "w", true, "use weak characters (2-9, A-N, P-Z, a-k, m-z)")
	rootCmd.Flags().StringP("exclude", "x", "", "excluded characters (will be removed)")
}

func resolveConfig(args []string) (arg.ConfigType, error) {
	cfg := arg.ConfigType{
		Symbols:   true,
		WeakChars: true,
	}
	fs := flag.NewFlagSet("CharVomit", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	fs.StringVar(&cfg.OutputFile, "output-file", "", "write the generated password to a file instead of stdout")
	fs.StringVar(&cfg.OutputFile, "o", "", "write the generated password to a file instead of stdout")
	fs.BoolVar(&cfg.Digits, "d", false, "use numeric digits")
	fs.BoolVar(&cfg.LowerCase, "l", false, "use lower-case letters")
	fs.BoolVar(&cfg.UpperCase, "u", false, "use upper-case letters")
	fs.BoolVar(&cfg.Symbols, "s", true, "use symbols: !#%+:=?@")
	fs.BoolVar(&cfg.Symbols, "symbols", true, "use symbols: !#%+:=?@")
	fs.BoolVar(&cfg.WeakChars, "w", true, "use weak characters (2-9, A-N, P-Z, a-k, m-z)")
	fs.StringVar(&cfg.Excluded, "x", "", "excluded characters (will be removed)")
	fs.BoolVar(&cfg.ShowHelp, "h", false, "show help and exit")
	fs.BoolVar(&cfg.Version, "v", false, "show version")

	if err := fs.Parse(args); err != nil {
		return cfg, err
	}

	if fs.NArg() > 1 {
		return cfg, fmt.Errorf("too many arguments: expected at most 1 positional length, got %d", fs.NArg())
	}

	cfg.PasswordLen = 32
	if fs.NArg() == 1 {
		parsedLen, err := strconv.Atoi(fs.Arg(0))
		if err != nil {
			return cfg, fmt.Errorf("cannot parse argument '%s': %w", fs.Arg(0), err)
		}

		cfg.PasswordLen = parsedLen
		if cfg.PasswordLen < 0 {
			cfg.PasswordLen *= -1
		}
	}

	return cfg, nil
}

func resolveCommandConfig(cmd *cobra.Command, args []string) (arg.ConfigType, error) {
	cfg, err := resolveConfig(args)
	if err != nil {
		return cfg, err
	}

	outputFile, err := cmd.Flags().GetString("output-file")
	if err != nil {
		return cfg, err
	}
	digits, err := cmd.Flags().GetBool("digits")
	if err != nil {
		return cfg, err
	}
	lowercase, err := cmd.Flags().GetBool("lowercase")
	if err != nil {
		return cfg, err
	}
	uppercase, err := cmd.Flags().GetBool("uppercase")
	if err != nil {
		return cfg, err
	}
	symbols, err := cmd.Flags().GetBool("symbols")
	if err != nil {
		return cfg, err
	}
	weak, err := cmd.Flags().GetBool("weak")
	if err != nil {
		return cfg, err
	}
	excluded, err := cmd.Flags().GetString("exclude")
	if err != nil {
		return cfg, err
	}

	cfg.OutputFile = outputFile
	cfg.Digits = digits
	cfg.LowerCase = lowercase
	cfg.UpperCase = uppercase
	cfg.Symbols = symbols
	cfg.WeakChars = weak
	cfg.Excluded = excluded

	return cfg, nil
}

func writePassword(password, outputFile string, _ bool) error {
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(password+"\n"), 0o600); err != nil {
			return err
		}
		return nil
	}

	if outputFile == "" {
		fmt.Println(password)
	}

	return nil
}

func run(cfg arg.ConfigType) error {
	var cv CharVomit.CharVomit
	if err := cv.SetAcceptableChars(cfg); err != nil {
		return fmt.Errorf("could not set acceptable chars: %w", err)
	}

	pw, err := cv.Puke(cfg.PasswordLen)
	if err != nil {
		return fmt.Errorf("Puke(%d) error: %w\n", cfg.PasswordLen, err)
	}

	return writePassword(pw, cfg.OutputFile, cfg.Stdout)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
