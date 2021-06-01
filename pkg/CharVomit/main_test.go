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
				if !strings.ContainsRune(s, rune(c)) {
					t.Errorf("char %q not present in %s", c, s)
				}
			}
		}

		// This should return error
		_, err := cv.Puke(0)
		if err == nil {
			t.Errorf("err was nil for Puke(0)")
		}
	}
}

