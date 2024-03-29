package CharVomit

import (
	"strings"
	"testing"

	"github.com/mojotx/CharVomit/pkg/arg"
	"github.com/stretchr/testify/assert"
)

func TestCharVomit_Puke(t *testing.T) {

	charSets := []string{
		AllUpperCase,
		AllLowerCase,
		DefaultSymbols,
		DefaultChars,
	}

	for _, s := range charSets {

		cv := NewCharVomit(s)

		for i := 1; i < len(s); i++ {
			pw, err := cv.Puke(i)

			if err != nil {
				t.Errorf("iteration i=%d err=%s", i, err.Error())
			}
			//assert.Nilf(t, err, "iteration i=%d err=%s", i, err.Error())

			if len(pw) != i {
				t.Errorf("i=%d: len(pw) = %d", i, len(pw))
			}

			for _, c := range pw {
				if !strings.ContainsRune(s, c) {
					t.Errorf("char %q not present in %s", c, s)
				}
			}
		}

		// This should return error
		_, err := cv.Puke(0)
		assert.NotNil(t, err, "err was nil for Puke(0)")

		// Try negative
		_, err = cv.Puke(-1)
		assert.NotNil(t, err, "err was nil for Puke(-1)")
	}
}

func TestUC(t *testing.T) {
	cv := NewCharVomit(AllUpperCase)

	pw, err := cv.Puke(99)
	if err != nil {
		t.Errorf("Received error upper-case puker: %s", err.Error())
	}

	for i, c := range pw {
		if !strings.ContainsRune(AllUpperCase, c) {
			t.Errorf("invalid char '%c' at offset %d of string '%s'", c, i, pw)
		}
	}
}

func TestLC(t *testing.T) {
	cv := NewCharVomit(AllLowerCase)

	pw, err := cv.Puke(99)
	if err != nil {
		t.Errorf("Received error lower-case puker: %s", err.Error())
	}

	for i, c := range pw {
		if !strings.ContainsRune(AllLowerCase, c) {
			t.Errorf("invalid char '%c' at offset %d of string '%s'", c, i, pw)
		}
	}
}
func TestDigits(t *testing.T) {
	cv := NewCharVomit(AllDigits)

	pw, err := cv.Puke(99)
	if err != nil {
		t.Errorf("Received error digit puker: %s", err.Error())
	}

	for i, c := range pw {
		if !strings.ContainsRune(AllDigits, c) {
			t.Errorf("invalid char '%c' at offset %d of string '%s'", c, i, pw)
		}
	}
}

func TestWeakChars(t *testing.T) {
	cv := NewCharVomit(WeakChars)

	pw, err := cv.Puke(99)
	if err != nil {
		t.Errorf("Received error weak puker: %s", err.Error())
	}

	for i, c := range pw {
		if !strings.ContainsRune(WeakChars, c) {
			t.Errorf("invalid char '%c' at offset %d of string '%s'", c, i, pw)
		}
	}
}

func TestDefaults(t *testing.T) {
	cv := NewCharVomit(DefaultChars)

	pw, err := cv.Puke(99)
	if err != nil {
		t.Errorf("Received error default puker: %s", err.Error())
	}

	for i, c := range pw {
		if !strings.ContainsRune(DefaultChars, c) {
			t.Errorf("invalid char '%c' at offset %d of string '%s'", c, i, pw)
		}
	}
}

func TestRemoveExcluded(t *testing.T) {
	// Set up a CharVomit instance with a predefined set of acceptable characters.
	cv := NewCharVomit(DefaultChars)

	// Define a configuration that specifies characters to exclude.
	config := arg.ConfigType{
		Excluded: "0O1l",
	}

	// Call RemoveExcluded with the configuration.
	err := cv.RemoveExcluded(config)

	// Check if an error occurred.
	if err != nil {
		t.Errorf("RemoveExcluded() returned an error: %v", err)
	}

	// Check if the excluded characters have been removed from AcceptableChars.
	// cSpell:disable-next-line
	expectedChars := "!#%+23456789:=?@ABCDEFGHJKLMNPRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	if cv.AcceptableChars != expectedChars {
		t.Errorf("RemoveExcluded() did not remove the expected characters. Expected: %s, Got: %s", expectedChars, cv.AcceptableChars)
	}

	// Test with an empty Excluded string to ensure no error is returned.
	config.Excluded = ""
	err = cv.RemoveExcluded(config)
	if err != nil {
		t.Errorf("RemoveExcluded() returned an error when Excluded was empty: %s", err.Error())
	}

	// Test with all characters excluded to ensure an error is returned.
	config.Excluded = DefaultChars
	err = cv.RemoveExcluded(config)
	assert.NotNil(t, err, "RemoveExcluded() did not return an error when all characters were excluded")
}

func TestSetAcceptableChars(t *testing.T) {
	t.Run("Default configuration", func(t *testing.T) {
		cv := NewCharVomit("")
		config := arg.ConfigType{}
		err := cv.SetAcceptableChars(config)
		assert.NoError(t, err)
		assert.Equal(t, DefaultChars, cv.AcceptableChars)
	})

	t.Run("Weak characters only", func(t *testing.T) {
		cv := NewCharVomit("")
		config := arg.ConfigType{
			WeakChars: true,
		}
		err := cv.SetAcceptableChars(config)
		assert.NoError(t, err)
		assert.Equal(t, WeakChars, cv.AcceptableChars)
	})

	t.Run("Digits, UpperCase, and LowerCase", func(t *testing.T) {
		cv := NewCharVomit("")
		config := arg.ConfigType{
			Digits:    true,
			UpperCase: true,
			LowerCase: true,
		}
		err := cv.SetAcceptableChars(config)
		assert.NoError(t, err)
		expectedChars := AllDigits + AllUpperCase + AllLowerCase
		assert.Equal(t, expectedChars, cv.AcceptableChars)
	})

	t.Run("Symbols only", func(t *testing.T) {
		cv := NewCharVomit("")
		config := arg.ConfigType{
			Symbols: true,
		}
		err := cv.SetAcceptableChars(config)
		assert.NoError(t, err)
		assert.Equal(t, DefaultSymbols, cv.AcceptableChars)
	})

	t.Run("Weak characters with other flags", func(t *testing.T) {
		cv := NewCharVomit("")
		config := arg.ConfigType{
			WeakChars: true,
			Digits:    true,
		}
		err := cv.SetAcceptableChars(config)
		assert.Error(t, err)
	})

	t.Run("Excluded characters", func(t *testing.T) {
		cv := NewCharVomit("")
		config := arg.ConfigType{
			Excluded: "0O1l",
		}
		err := cv.SetAcceptableChars(config)
		assert.NoError(t, err)
		// cSpell:disable-next-line
		expectedChars := "!#%+23456789:=?@ABCDEFGHJKLMNPRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
		assert.Equal(t, expectedChars, cv.AcceptableChars)
	})
}
