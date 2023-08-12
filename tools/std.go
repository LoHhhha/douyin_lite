package tools

import (
	"fmt"
)

func ErrorPrint(err error) {
	fmt.Printf("[Error] %v.\n", err)
}
