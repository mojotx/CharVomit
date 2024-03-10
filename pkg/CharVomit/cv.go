package CharVomit

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/mojotx/CharVomit/pkg/arg"
)

// CharVomit package generates random passwords
// uc, lc, d, s, w
const (
	AllUpperCase   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AllLowerCase   = "abcdefghijklmnopqrstuvwxyz"
	AllDigits      = "0123456789"
	DefaultSymbols = "!#%+:=?@"
	// cSpell:disable-next-line
	DefaultChars = "!#%+23456789:=?@ABCDEFGHJKLMNPRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	WeakChars    = "23456789ABCDEFGHJKLMNPRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

type CharVomit struct {
	AcceptableChars string
}

func NewCharVomit(acceptableChars string) CharVomit {
	return CharVomit{
		AcceptableChars: acceptableChars,
	}
}

func (c *CharVomit) getByte(i int64) byte {
	return c.AcceptableChars[i]
}

// acLen is a helper function that returns the length of the acceptable chars slice
// It is intended to be used by the Puke() function
func (c *CharVomit) acLen() *big.Int {
	return big.NewInt(int64(len(c.AcceptableChars)))
}

// Puke returns a string of characters of a given length
func (c *CharVomit) Puke(length int) (string, error) {
	var b strings.Builder

	// Check for invalid length <= 0
	if length <= 0 {
		return b.String(), fmt.Errorf("password length %d <= 0", length)
	}

	for i := 0; i < length; i++ {
		ri, err := rand.Int(rand.Reader, c.acLen())
		if err != nil {
			return "", err
		}

		b.WriteByte(c.getByte(ri.Int64()))
	}

	return b.String(), nil
}

func (c *CharVomit) RemoveExcluded(config arg.ConfigType) error {
	if len(config.Excluded) <= 0 || len(c.AcceptableChars) <= 0 {
		return nil
	}

	for _, r := range config.Excluded {
		c.AcceptableChars = strings.ReplaceAll(c.AcceptableChars, string(r), "")
	}

	if len(c.AcceptableChars) <= 0 {
		return fmt.Errorf("no acceptable characters")
	}
	return nil
}

func (c *CharVomit) SetAcceptableChars(config arg.ConfigType) error {
	// Initialize to nothing
	c.AcceptableChars = ""

	// Weak characters cannot be used with other flags
	if config.WeakChars {
		if config.Symbols {
			return fmt.Errorf("cannot specify weak characters with symbols")
		}

		if config.Digits || config.UpperCase || config.LowerCase {
			return fmt.Errorf("redundant specification of characters with weak characters")
		}

		c.AcceptableChars = WeakChars
		return nil
	}

	if config.Digits {
		c.AcceptableChars += AllDigits
	}

	if config.UpperCase {
		c.AcceptableChars += AllUpperCase
	}

	if config.LowerCase {
		c.AcceptableChars += AllLowerCase
	}

	if config.Symbols {
		c.AcceptableChars += DefaultSymbols
	}

	if len(c.AcceptableChars) <= 0 {
		c.AcceptableChars = DefaultChars
	}

	if err := c.RemoveExcluded(config); err != nil {
		return err
	}

	return nil
}
