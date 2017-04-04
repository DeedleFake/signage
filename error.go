package signage

import (
	"fmt"
	"runtime"
)

type errorImpl struct {
	error

	file string
	line int
}

func Error(err error) error {
	_, file, line, ok := runtime.Caller(0)
	if !ok {
		panic(fmt.Errorf("Failed to get caller frame for error: %q", err))
	}

	return errorImpl{
		error: err,

		file: file,
		line: line,
	}
}

func (err errorImpl) Error() string {
	return fmt.Sprintf("%v at %v:%v", err.error, err.file, err.line)
}
