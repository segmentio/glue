package log

import "fmt"

var DebugMode bool

func Print(v ...interface{}) {
	fmt.Println(v...)
}

func Printf(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

func Debug(v ...interface{}) {
	if DebugMode {
		Print(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if DebugMode {
		Printf(format, v...)
	}
}
