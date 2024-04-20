package util

import (
	"fmt"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// FatalIf exits if the error is not nil
func FatalIf(err error) {
	if err != nil {
		if stackErr, ok := err.(stackTracer); ok {
			logrus.WithField("stacktrace", fmt.Sprintf("%+v", stackErr.StackTrace()))
		} else {
			debug.PrintStack()
		}
		logrus.Fatalf("Fatal error: %s\n", err)
	}
}
