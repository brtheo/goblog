package util

import "fmt"

func Throw(err error) {
	if err != nil {
		fmt.Print(err)
	}
}
