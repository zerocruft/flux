package debug

import (
	"fmt"
)

var (
	isDebug bool
)

func Init(is bool) {
	isDebug = is
}

func Log(msg interface{}) {
	if isDebug {
		fmt.Print("DEBUG: ")
		fmt.Println(msg)
	}
}
