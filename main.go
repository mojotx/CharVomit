package main

import (
	"fmt"
	"github.com/mojotx/CharVomit/pkg/CharVomit"
	"os"
	"strconv"
)

// TO-DO:
// - Add support for duplicate character checking
// - Add support for specifying valid characters, e.g., all upper-case, etc.
func main() {
	cv := CharVomit.NewCharVomit(CharVomit.DefaultChars)

	pwLen := 32

	if len(os.Args) == 2 {

		var err error
		pwLen, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if pwLen < 0 {
			pwLen = pwLen * -1
		}
	}
	pw, err := cv.Puke(pwLen)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(pw)
}
