package gee

import (
	"fmt"
	"runtime"
)

func Tract() {
	var ptr [32]uintptr
	n := runtime.Callers(3, ptr[:])
	for _, pc := range ptr[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		fmt.Println(file, line)
	}
}

func Recoverier() HandleFunc {
	return func(c *Context) {
		defer func() {
			Tract()
			if err := recover(); err != nil {

			}
		}()

		c.Next()
	}
}
