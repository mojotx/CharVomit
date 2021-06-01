package CharVomit

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// CharVomit package generates random passwords
const (
	AllUpperCase   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AllLowerCase   = "abcdefghijklmnopqrstuvwxyz"
	AllDigits      = "0123456789"
	DefaultSymbols = "!#%+:=?@"
	DefaultChars   = "!#%+23456789:=?@ABCDEFGHJKLMNPRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	WeakChars      = "23456789ABCDEFGHJKLMNPRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

type CharVomit struct {
	acceptableChars string
}

func NewCharVomit(s string) CharVomit {
	return CharVomit{
		acceptableChars: s,
	}
}

func (c *CharVomit) getByte(i int64) byte {
	return c.acceptableChars[i]
}

// acLen is a helper function that returns the length of the acceptable chars slice
// It is intended to be used by the Puke() function
func (c *CharVomit) acLen() *big.Int {
	return big.NewInt(int64(len(c.acceptableChars)))
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
