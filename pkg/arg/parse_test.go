package arg

import (
	"bytes"
	"flag"
	"os"
	"regexp"
	"strings"
	"testing"
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

	if !re.MatchString(haystack) {
		t.Errorf("Did not find '%s' in '%s'", needle, haystack)
	}
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

	if Config.UpperCase {
		t.Error("Config.UpperCase should not be set")
	}

	if Config.LowerCase {
		t.Error("Config.LowerCase should not be set")
	}

	if Config.Digits {
		t.Error("Config.Digits should not be set")
	}

	if Config.Symbols {
		t.Error("Config.Symbols should not be set")
	}

	if Config.ShowHelp {
		t.Error("Config.ShowHelp should not be set")
	}

	if Config.Version {
		t.Error("Config.Version should not be set")
	}

	if Config.Excluded != "" {
		t.Errorf("Config.Excluded should be empty (%s)", Config.Excluded)
	}

	if !Config.WeakChars {
		t.Error("Config.WeakChars *SHOULD* be set")
	}
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

	if Config.UpperCase {
		t.Error("Config.UpperCase should not be set")
	}

	if Config.LowerCase {
		t.Error("Config.LowerCase should not be set")
	}

	if Config.Digits {
		t.Error("Config.Digits should not be set")
	}

	if Config.Symbols {
		t.Error("Config.Symbols should not be set")
	}

	if Config.ShowHelp {
		t.Error("Config.ShowHelp should not be set")
	}

	if Config.Version {
		t.Error("Config.Version should not be set")
	}

	if Config.Excluded != "" {
		t.Errorf("Config.Excluded should be empty (%s)", Config.Excluded)
	}

	if Config.WeakChars {
		t.Error("Config.WeakChars should not be set")
	}

}

func TestParseVersion(t *testing.T) {
	// idempotency for the win
	originalWriter := flag.CommandLine.Output()
	defer func() {
		flag.CommandLine.SetOutput(originalWriter)
	}()

	// Intercept output
	var buf bytes.Buffer
	flag.CommandLine.SetOutput(&buf)

	Config = ConfigType{}
	newArgs := []string{os.Args[0], "-v"}
	savedArgs := os.Args
	defer func() {
		os.Args = savedArgs
	}()

	os.Args = newArgs

	fs := flag.NewFlagSet("TestParseNoArg", flag.ExitOnError)
	Parse(fs)

	if Config.UpperCase {
		t.Error("Config.UpperCase should not be set")
	}

	if Config.LowerCase {
		t.Error("Config.LowerCase should not be set")
	}

	if Config.Digits {
		t.Error("Config.Digits should not be set")
	}

	if Config.Symbols {
		t.Error("Config.Symbols should not be set")
	}

	if Config.ShowHelp {
		t.Error("Config.ShowHelp should not be set")
	}

	if !Config.Version {
		t.Error("Config.Version SHOULD be set")
	}

	if Config.Excluded != "" {
		t.Errorf("Config.Excluded should be empty (%s)", Config.Excluded)
	}

	if Config.WeakChars {
		t.Error("Config.WeakChars should not be set")
	}

	version := strings.TrimRight(buf.String(), "\r\n\t ")

	if version != Version() {
		t.Errorf("invalid version, wanted '%s' got '%s'", Version(), version)
	}
}
