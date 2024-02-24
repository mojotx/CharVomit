package CharVomit

import (
	"strings"
	"testing"
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
		if err == nil {
			t.Errorf("err was nil for Puke(0)")
		}

		// Try negative
		_, err = cv.Puke(-1)
		if err == nil {
			t.Errorf("err was nil for Puke(-1)")
		}
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
		t.Errorf("Received error digit puker: %s", err.Error())
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
		t.Errorf("Received error digit puker: %s", err.Error())
	}

	for i, c := range pw {
		if !strings.ContainsRune(DefaultChars, c) {
			t.Errorf("invalid char '%c' at offset %d of string '%s'", c, i, pw)
		}
	}
}
