package arg

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const needle = `If a password length is not specified, 32 is used.`

// TestUsage tries to test the Usage() function, but it's hard to do
// since Go unit tests are a different binary. Therefore I am just
// going to look for the line, `If a password length is not specified,
// 32 is used.`
func TestUsage(t *testing.T) {
	// idempotency for the win
	originalWriter := flag.CommandLine.Output()
	defer func() {
		flag.CommandLine.SetOutput(originalWriter)
	}()

	// Set the flag output to a byte buffer
	var buf bytes.Buffer
	flag.CommandLine.SetOutput(&buf)
	Usage()

	haystack := buf.String()
	re := regexp.MustCompile(needle)

	assert.Regexp(t, re, haystack)
}

func TestParseWeak(t *testing.T) {

	Config = ConfigType{}
	newArgs := []string{os.Args[0], "-w"}
	savedArgs := os.Args
	defer func() {
		os.Args = savedArgs
	}()

	os.Args = newArgs

	fs := flag.NewFlagSet("TestParseWeak", flag.ExitOnError)
	Parse(fs)

	t.Logf("Config: %+v", Config)

	assert.False(t, Config.UpperCase)
	assert.False(t, Config.LowerCase)
	assert.False(t, Config.Digits)
	assert.False(t, Config.Symbols)
	assert.False(t, Config.ShowHelp)
	assert.False(t, Config.Version)
	assert.Empty(t, Config.Excluded)
	assert.True(t, Config.WeakChars)
}

func TestParseNoArg(t *testing.T) {
	// Now try no args
	Config = ConfigType{}
	newArgs := []string{os.Args[0]}
	savedArgs := os.Args
	defer func() {
		os.Args = savedArgs
	}()

	// Set args
	os.Args = newArgs

	fs := flag.NewFlagSet("TestParseNoArg", flag.ExitOnError)

	Parse(fs)

	assert.False(t, Config.UpperCase)
	assert.False(t, Config.LowerCase)
	assert.False(t, Config.Digits)
	assert.False(t, Config.Symbols)
	assert.False(t, Config.ShowHelp)
	assert.False(t, Config.Version)
	assert.Empty(t, Config.Excluded)
	assert.False(t, Config.WeakChars)

}

func TestParseVersion(t *testing.T) {
	Config = ConfigType{}
	newArgs := []string{os.Args[0], "-v"}
	savedArgs := os.Args
	defer func() {
		os.Args = savedArgs
	}()

	os.Args = newArgs

	fs := flag.NewFlagSet("TestParseVersion", flag.ExitOnError)

	// Intercept output
	var buf bytes.Buffer
	fs.SetOutput(&buf)

	exitAfter, rc := Parse(fs)
	t.Logf("exitAfter: %t, rc: %d", exitAfter, rc)

	assert.False(t, Config.UpperCase)
	assert.False(t, Config.LowerCase)
	assert.False(t, Config.Digits)
	assert.False(t, Config.Symbols)
	assert.False(t, Config.ShowHelp)
	assert.True(t, Config.Version)
	assert.Empty(t, Config.Excluded)
	assert.False(t, Config.WeakChars)

	actualVersion := strings.TrimRight(buf.String(), "\r\n\t ")

	expectedVersion := Version()
	t.Logf("Expected version: %s", expectedVersion)
	t.Logf("Actual version: %s", actualVersion)
	assert.Equal(t, expectedVersion, actualVersion)
}

func TestParse(t *testing.T) {
	// Save the original command-line arguments and restore them after the test.
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set up a mock command-line input.
	os.Args = []string{"CharVomit", "-u", "-l", "-d", "-s", "-w", "-x", "0O1l", "8"}

	// Create a new FlagSet to simulate command-line arguments.
	fs := flag.NewFlagSet("CharVomit", flag.ExitOnError)

	// suppress the help output
	oldOutput := fs.Output()
	var buf bytes.Buffer
	fs.SetOutput(&buf)

	// Call the Parse function with the mock FlagSet.
	exitAfter, rc := Parse(fs)

	// Check if the function exits as expected.
	assert.False(t, exitAfter)

	// Check if the return code is as expected.
	assert.Zero(t, rc)

	// Check if the parsed configuration matches the expected values.
	expectedConfig := ConfigType{
		PasswordLen: 8,
		Digits:      true,
		ShowHelp:    false,
		LowerCase:   true,
		Symbols:     true,
		UpperCase:   true,
		WeakChars:   true,
		Version:     false,
		Excluded:    "0O1l",
	}

	assert.Equal(t, expectedConfig, Config)
	fs.SetOutput(oldOutput)
}

func TestParseConfigReturnsValue(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	Config = ConfigType{}
	os.Args = []string{"CharVomit", "-u", "-l", "8"}

	fs := flag.NewFlagSet("CharVomit", flag.ContinueOnError)
	cfg, exitAfter, rc := ParseConfig(fs)

	assert.False(t, exitAfter)
	assert.Zero(t, rc)
	assert.Equal(t, 8, cfg.PasswordLen)
	assert.True(t, cfg.UpperCase)
	assert.True(t, cfg.LowerCase)
	assert.False(t, cfg.WeakChars)
	assert.Equal(t, Config, cfg)
}

func TestParseFlagSetCanBeReused(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"CharVomit", "-u", "8"}
	fs := flag.NewFlagSet("CharVomit", flag.ContinueOnError)
	_, rc := Parse(fs)
	assert.Zero(t, rc)
	assert.Equal(t, 8, Config.PasswordLen)
	assert.True(t, Config.UpperCase)
	assert.False(t, Config.LowerCase)

	os.Args = []string{"CharVomit", "-l", "12"}
	_, rc = Parse(fs)
	assert.Zero(t, rc)
	assert.Equal(t, 12, Config.PasswordLen)
	assert.False(t, Config.UpperCase)
	assert.True(t, Config.LowerCase)
}

func TestParseArgsReturnsValuesFromProvidedFlagSet(t *testing.T) {
	Config = ConfigType{PasswordLen: 99}
	localCfg := ConfigType{}
	fs := flag.NewFlagSet("CharVomit", flag.ContinueOnError)
	fs.BoolVar(&localCfg.UpperCase, "u", false, "use upper-case letters")
	fs.BoolVar(&localCfg.LowerCase, "l", false, "use lower-case letters")
	fs.BoolVar(&localCfg.Digits, "d", false, "use numeric digits")
	fs.BoolVar(&localCfg.Symbols, "s", false, "use symbols: !#%+:=?@")
	fs.BoolVar(&localCfg.WeakChars, "w", false, "use weak characters")
	fs.BoolVar(&localCfg.ShowHelp, "h", false, "show help and exit")
	fs.BoolVar(&localCfg.Version, "v", false, "show version")
	fs.StringVar(&localCfg.Excluded, "x", "", "excluded characters")

	cfg, exitAfter, rc := ParseArgs([]string{"-u", "-s", "12"}, fs)

	assert.False(t, exitAfter)
	assert.Zero(t, rc)
	assert.Equal(t, 12, cfg.PasswordLen)
	assert.True(t, cfg.UpperCase)
	assert.True(t, cfg.Symbols)
	assert.False(t, cfg.LowerCase)
	assert.Equal(t, ConfigType{PasswordLen: 99}, Config)
}

func TestParseArgsSupportsConcurrentFlagSets(t *testing.T) {
	type result struct {
		cfg ConfigType
		err error
	}

	args := [][]string{{"-u", "12"}, {"-l", "8"}}
	results := make(chan result, len(args))
	start := make(chan struct{})
	var wg sync.WaitGroup

	for _, testArgs := range args {
		wg.Add(1)
		go func(args []string) {
			defer wg.Done()
			fs := flag.NewFlagSet("CharVomit", flag.ContinueOnError)
			fs.SetOutput(io.Discard)
			<-start
			cfg, exitAfter, rc := ParseArgs(args, fs)
			if exitAfter || rc != 0 {
				results <- result{err: fmt.Errorf("unexpected parse result: exitAfter=%t rc=%d", exitAfter, rc)}
				return
			}
			results <- result{cfg: cfg}
		}(testArgs)
	}

	close(start)
	wg.Wait()
	close(results)

	var configs []ConfigType
	for result := range results {
		require.NoError(t, result.err)
		configs = append(configs, result.cfg)
	}

	require.Len(t, configs, 2)
	assert.ElementsMatch(t, []ConfigType{
		{PasswordLen: 12, UpperCase: true},
		{PasswordLen: 8, LowerCase: true},
	}, configs)
}

func TestParseArgsReportsCustomArgs(t *testing.T) {
	fs := flag.NewFlagSet("CharVomit", flag.ContinueOnError)
	cfg := ConfigType{}
	RegisterFlags(fs, &cfg)

	var buf bytes.Buffer
	fs.SetOutput(&buf)

	_, exitAfter, rc := ParseArgs([]string{"--bad-flag"}, fs)

	assert.True(t, exitAfter)
	assert.Equal(t, 1, rc)
	assert.Contains(t, buf.String(), "cannot parse args")
}

func TestParseHelp(t *testing.T) {
	// Save the original command-line arguments and restore them after the test.
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Initialize Config
	Config = ConfigType{}

	// Set up a mock command-line input.
	os.Args = []string{"CharVomit", "-h"}

	// Create a new FlagSet to simulate command-line arguments.
	fs := flag.NewFlagSet("CharVomit", flag.ExitOnError)

	// Suppress the help output
	oldOutput := fs.Output()
	defer fs.SetOutput(oldOutput)
	var buf bytes.Buffer
	fs.SetOutput(&buf)

	// Call the Parse function with the mock FlagSet
	exitAfter, rc := Parse(fs)

	// Check if the function exits as expected
	assert.True(t, exitAfter)

	// Check if the return code is as expected
	assert.Zero(t, rc, "Expected non-zero return code for help output")

	// Check if the parsed configuration matches the expected values
	expectedConfig := ConfigType{
		PasswordLen: 32,
		Digits:      false,
		ShowHelp:    true,
		LowerCase:   false,
		Symbols:     false,
		UpperCase:   false,
		WeakChars:   false,
		Version:     false,
		Excluded:    "",
	}
	assert.Equal(t, expectedConfig, Config, "Parsed configuration does not match expected values")

	// t.Logf("Help output: %s", buf.String())

	// Check if the help output is empty
	t.Logf("Size of buf is %d", buf.Len())
	helpOutput := buf.String()
	assert.NotEmpty(t, helpOutput, "Expected help output to be non-empty")
}

func TestParseInvalidMultipleLengths(t *testing.T) {
	Config = ConfigType{}
	newArgs := []string{os.Args[0], "8", "9"}
	savedArgs := os.Args
	defer func() {
		os.Args = savedArgs
	}()

	os.Args = newArgs

	fs := flag.NewFlagSet("TestParseInvalidMultipleLengths", flag.ExitOnError)
	var buf bytes.Buffer
	fs.SetOutput(&buf)

	exitAfter, rc := Parse(fs)

	assert.True(t, exitAfter)
	assert.Equal(t, 1, rc)
	assert.Contains(t, buf.String(), "too many arguments")
}

func TestParseInvalidLength(t *testing.T) {
	Config = ConfigType{}
	newArgs := []string{os.Args[0], "not-a-number"}
	savedArgs := os.Args
	defer func() {
		os.Args = savedArgs
	}()

	os.Args = newArgs

	fs := flag.NewFlagSet("TestParseInvalidLength", flag.ExitOnError)
	var buf bytes.Buffer
	fs.SetOutput(&buf)

	exitAfter, rc := Parse(fs)

	assert.True(t, exitAfter)
	assert.Equal(t, 1, rc)
	assert.Contains(t, buf.String(), "cannot parse argument")
}
