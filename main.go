package main

import (
	"fmt"
	"github.com/mojotx/CharVomit/pkg/CharVomit"
)

func main() {
	cv := CharVomit.NewCharVomit(CharVomit.DefaultChars)

	pw, err := cv.Puke(32)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(pw)
}
