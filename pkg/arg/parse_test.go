package arg

import (
	"bytes"
	"flag"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

	t.Logf("Config: %+v", Config)

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
