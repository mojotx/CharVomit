package main

import (
	"strings"
	"testing"

	"github.com/mojotx/CharVomit/pkg/CharVomit"
)

func TestUC(t *testing.T) {
	cv := CharVomit.NewCharVomit(CharVomit.AllUpperCase)

	pw, err := cv.Puke(99)
	if err != nil {
		t.Errorf("Received error upper-case puker: %s", err.Error())
	}

	for i, c := range pw {
		if !strings.ContainsRune(CharVomit.AllUpperCase, c) {
			t.Errorf("invalid char '%c' at offset %d of string '%s'", rune(c), i, pw)
		}
	}
}

func TestLC(t *testing.T) {
	cv := CharVomit.NewCharVomit(CharVomit.AllLowerCase)

	pw, err := cv.Puke(99)
	if err != nil {
		t.Errorf("Received error lower-case puker: %s", err.Error())
	}

	for i, c := range pw {
		if !strings.ContainsRune(CharVomit.AllLowerCase, c) {
			t.Errorf("invalid char '%c' at offset %d of string '%s'", rune(c), i, pw)
		}
	}
}

func TestDigits(t *testing.T) {
	cv := CharVomit.NewCharVomit(CharVomit.AllDigits)

	pw, err := cv.Puke(99)
	if err != nil {
		t.Errorf("Received error digit puker: %s", err.Error())
	}

	for i, c := range pw {
		if !strings.ContainsRune(CharVomit.AllDigits, c) {
			t.Errorf("invalid char '%c' at offset %d of string '%s'", rune(c), i, pw)
		}
	}
}

func TestWeakChars(t *testing.T) {
	cv := CharVomit.NewCharVomit(CharVomit.WeakChars)

	pw, err := cv.Puke(99)
	if err != nil {
		t.Errorf("Received error digit puker: %s", err.Error())
	}

	for i, c := range pw {
		if !strings.ContainsRune(CharVomit.WeakChars, c) {
			t.Errorf("invalid char '%c' at offset %d of string '%s'", rune(c), i, pw)
		}
	}
}

func TestDefaults(t *testing.T) {
	cv := CharVomit.NewCharVomit(CharVomit.DefaultChars)

	pw, err := cv.Puke(99)
	if err != nil {
		t.Errorf("Received error digit puker: %s", err.Error())
	}

	for i, c := range pw {
		if !strings.ContainsRune(CharVomit.DefaultChars, c) {
			t.Errorf("invalid char '%c' at offset %d of string '%s'", rune(c), i, pw)
		}
	}
}
