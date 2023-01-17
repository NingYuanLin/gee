package gee

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func trace(message string) string {
	var pcs [32]uintptr
	// Skip: 0: runtime.Callers itself; 1: trace(); 2:Recovery();
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message)
	str.WriteString("\nTraceback:")

	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		str.WriteString(fmt.Sprintf("\n\t%s:%d", frame.File, frame.Line))
		if more == false {
			break
		}
	}
	return str.String()
}

func Recovery(f func(c *Context)) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				traceInfo := trace(message)
				log.Println(traceInfo)
				f(c)
			}
		}()
		c.Next()
	}
}
