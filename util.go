package main

import (
	"fmt"
	"runtime"
	"strings"
)

func errorf(message string, a ...interface{}) string {
	pc, _, _, _ := runtime.Caller(1)
	s := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	newMessage := fmt.Sprintf("[%s] %s", s[len(s)-1], message)
	return fmt.Sprintf(newMessage, a...)
}
