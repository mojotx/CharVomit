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

func TestParse(t *testing.T) {
	// Save the original command-line arguments and restore them after the test.
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set up a mock command-line input.
	os.Args = []string{"CharVomit", "-u", "-l", "-d", "-s", "-w", "-x", "0O1l", "8"}

	// Create a new FlagSet to simulate command-line arguments.
	fs := flag.NewFlagSet("CharVomit", flag.ExitOnError)

	// Call the Parse function with the mock FlagSet.
	exitAfter, rc := Parse(fs)

	// Check if the function exits as expected.
	if exitAfter {
		t.Errorf("Parse() exited unexpectedly")
	}

	// Check if the return code is as expected.
	if rc != 0 {
		t.Errorf("Parse() returned unexpected return code: %d", rc)
	}

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

	if Config != expectedConfig {
		t.Errorf("Parse() parsed configuration does not match expected values, wanted '%+v' got '%+v'", expectedConfig, Config)
	}
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
	if !exitAfter {
		t.Errorf("Parse() did not exit as expected")
	}

	// Check if the return code is as expected
	if rc != 0 {
		t.Errorf("Parse() returned unexpected return code: %d", rc)
	}

	// Check if the parsed configuration matches the expected values
	expectedConfig := ConfigType{
		PasswordLen: 0,
		Digits:      false,
		ShowHelp:    true,
		LowerCase:   false,
		Symbols:     false,
		UpperCase:   false,
		WeakChars:   false,
		Version:     false,
		Excluded:    "",
	}
	if Config != expectedConfig {
		t.Errorf("Parse() parsed configuration does not match expected values, wanted '%+v' got '%+v'", expectedConfig, Config)
	}

	// t.Logf("Help output: %s", buf.String())

	// Check if the help output is empty
	helpOutput := buf.String()
	if helpOutput != "" {
		t.Errorf("Parse() produced unexpected help output: %s", helpOutput)
	}
}
