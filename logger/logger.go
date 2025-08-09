package logger

import (
	"fmt"
	"os"
)

func Fail(a ...any) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func FailF(format string, a ...any) {
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(format, a...))
	os.Exit(1)
}
