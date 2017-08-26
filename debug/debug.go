package debug

import (
	"fmt"
)

var (
	isDebug func() bool
)

func Init(fn func() bool) {
	isDebug = fn
}



func Log(msg interface{}) {
	if isDebug() {
		fmt.Print("DEBUG: ")
		fmt.Println(msg)
	}
}
