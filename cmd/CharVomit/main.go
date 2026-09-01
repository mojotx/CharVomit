package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/mojotx/CharVomit/v2/pkg/CharVomit"
	"github.com/mojotx/CharVomit/v2/pkg/arg"
)

var rootCmd = &cobra.Command{
	Use:           "CharVomit [length]",
	Short:         "Generate random passwords.",
	Version:       arg.Version(),
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
	rootCmd.Flags().BoolP("symbols", "s", false, "use symbols: !#%+:=?@")
	rootCmd.Flags().BoolP("weak", "w", false, "use weak characters (2-9, A-N, P-Z, a-k, m-z)")
	rootCmd.Flags().StringP("exclude", "x", "", "excluded characters (will be removed)")
}

// resolveConfig parses the optional positional password length. rootCmd's
// Args: cobra.MaximumNArgs(1) already guarantees len(args) <= 1 here.
func resolveConfig(args []string) (arg.ConfigType, error) {
	cfg := arg.ConfigType{PasswordLen: 32}
	if len(args) == 0 {
		return cfg, nil
	}

	parsedLen, err := strconv.Atoi(args[0])
	if err != nil {
		return cfg, fmt.Errorf("cannot parse argument '%s': %w", args[0], err)
	}

	if parsedLen < 0 {
		negated := -parsedLen
		if negated < 0 {
			return cfg, fmt.Errorf("password length %d is out of range", parsedLen)
		}
		parsedLen = negated
	}
	cfg.PasswordLen = parsedLen

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

func writePassword(password, outputFile string, output io.Writer) error {
	if outputFile != "" {
		file, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			return err
		}
		if err := file.Chmod(0o600); err != nil {
			return errors.Join(err, file.Close())
		}
		if _, err := fmt.Fprintln(file, password); err != nil {
			return errors.Join(err, file.Close())
		}
		return file.Close()
	}

	_, err := fmt.Fprintln(output, password)
	return err
}

func run(cfg arg.ConfigType) error {
	var cv CharVomit.CharVomit
	if err := cv.SetAcceptableChars(cfg); err != nil {
		return fmt.Errorf("could not set acceptable chars: %w", err)
	}

	pw, err := cv.Puke(cfg.PasswordLen)
	if err != nil {
		return fmt.Errorf("Puke(%d) error: %w", cfg.PasswordLen, err)
	}

	return writePassword(pw, cfg.OutputFile, os.Stdout)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
