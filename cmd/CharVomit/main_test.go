package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mojotx/CharVomit/pkg/CharVomit"
	"github.com/mojotx/CharVomit/pkg/arg"
)

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func TestResolveConfigUses32CharacterDefault(t *testing.T) {
	cfg, err := resolveConfig(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.PasswordLen != 32 {
		t.Fatalf("expected default password length 32, got %d", cfg.PasswordLen)
	}
	if cfg.WeakChars || cfg.Symbols {
		t.Fatalf("expected no character class flags by default, got %+v", cfg)
	}

	var cv CharVomit.CharVomit
	if err := cv.SetAcceptableChars(cfg); err != nil {
		t.Fatalf("unexpected error setting default characters: %v", err)
	}

	expectedPool := CharVomit.DefaultChars
	if cv.AcceptableChars != expectedPool {
		t.Fatalf("expected default character pool %q, got %q", expectedPool, cv.AcceptableChars)
	}

	if strings.ContainsAny(cv.AcceptableChars, "0O1l") {
		t.Fatal("default character pool contains ambiguous characters")
	}
}

func TestResolveConfigAcceptsSymbolsAlias(t *testing.T) {
	cfg, err := resolveConfig([]string{"-s", "-u", "12"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.PasswordLen != 12 {
		t.Fatalf("expected password length 12, got %d", cfg.PasswordLen)
	}

	if !cfg.UpperCase {
		t.Fatal("expected uppercase characters to be enabled")
	}

	if !cfg.Symbols {
		t.Fatal("expected symbols to be enabled")
	}
}

func TestWritePasswordToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "password.txt")

	if err := writePassword("abc123", path, nil); err != nil {
		t.Fatalf("unexpected error writing file: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}

	if got := string(content); got != "abc123\n" {
		t.Fatalf("unexpected file content: %q", got)
	}
}

func TestWritePasswordReturnsStandardOutputError(t *testing.T) {
	want := errors.New("broken pipe")
	err := writePassword("abc123", "", failingWriter{err: want})
	if !errors.Is(err, want) {
		t.Fatalf("expected stdout error %v, got %v", want, err)
	}
}

func TestRootCommandPrintsVersion(t *testing.T) {
	var output bytes.Buffer
	rootCmd.SetArgs([]string{"--version"})
	rootCmd.SetOut(&output)
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(nil)
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := output.String(), "CharVomit version "+arg.Version()+"\n"; got != want {
		t.Fatalf("expected version output %q, got %q", want, got)
	}
}
