package main

import (
	"bytes"
	"errors"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/mojotx/CharVomit/v2/pkg/CharVomit"
	"github.com/mojotx/CharVomit/v2/pkg/arg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func TestResolveConfigUses32CharacterDefault(t *testing.T) {
	cfg, err := resolveConfig(nil)
	require.NoError(t, err)

	assert.Equal(t, 32, cfg.PasswordLen)
	assert.False(t, cfg.WeakChars)
	assert.False(t, cfg.Symbols)

	var cv CharVomit.CharVomit
	require.NoError(t, cv.SetAcceptableChars(cfg))

	expectedPool := CharVomit.DefaultChars
	assert.Equal(t, expectedPool, cv.AcceptableChars)
	assert.False(t, strings.ContainsAny(cv.AcceptableChars, "0O1l"))
}

func TestResolveConfigRejectsMinIntLength(t *testing.T) {
	// math.MinInt negated overflows back to itself; make sure that's reported
	// as an out-of-range error instead of silently passing through negative.
	_, err := resolveConfig([]string{strconv.Itoa(math.MinInt)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "out of range")
}

func TestResolveCommandConfigAcceptsSymbolsAlias(t *testing.T) {
	require.NoError(t, rootCmd.ParseFlags([]string{"-s", "-u"}))
	t.Cleanup(func() {
		require.NoError(t, rootCmd.Flags().Set("symbols", "false"))
		require.NoError(t, rootCmd.Flags().Set("uppercase", "false"))
	})

	cfg, err := resolveCommandConfig(rootCmd, []string{"12"})
	require.NoError(t, err)

	assert.Equal(t, 12, cfg.PasswordLen)
	assert.True(t, cfg.UpperCase)
	assert.True(t, cfg.Symbols)
}

func TestWritePasswordToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "password.txt")

	require.NoError(t, writePassword("abc123", path, nil))

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "abc123\n", string(content))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestWritePasswordRestrictsExistingFilePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "password.txt")
	require.NoError(t, os.WriteFile(path, []byte("old password\n"), 0o644))
	require.NoError(t, os.Chmod(path, 0o644))

	require.NoError(t, writePassword("abc123", path, nil))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestWritePasswordReturnsStandardOutputError(t *testing.T) {
	want := errors.New("broken pipe")
	err := writePassword("abc123", "", failingWriter{err: want})
	assert.ErrorIs(t, err, want)
}

func TestRootCommandPrintsVersion(t *testing.T) {
	var output bytes.Buffer
	rootCmd.SetArgs([]string{"--version"})
	rootCmd.SetOut(&output)
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(nil)
	})

	require.NoError(t, rootCmd.Execute())

	assert.Equal(t, "CharVomit version "+arg.Version()+"\n", output.String())
}
