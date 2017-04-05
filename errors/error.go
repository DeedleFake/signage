package errors

import (
	"fmt"
	"runtime"
)

type lineError struct {
	error

	file string
	line int
}

// Err returns an error which wraps the given error and provides the
// filename and line from which it was called.
func Err(err error) error {
	if err == nil {
		return nil
	}

	_, file, line, ok := runtime.Caller(0)
	if !ok {
		panic(fmt.Errorf("Failed to get caller frame for error: %q", err))
	}

	return lineError{
		error: err,

		file: file,
		line: line,
	}
}

func (err lineError) Error() string {
	return fmt.Sprintf("%v at %v:%v", err.error, err.file, err.line)
}
