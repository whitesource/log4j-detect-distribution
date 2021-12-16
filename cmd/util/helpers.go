package util

import (
	"fmt"
	"os"
	"strings"
)

const DefaultErrorExitCode = 1

var fatalErrHandler = fatal

// BehaviorOnFatal allows you to override the default behavior when a fatal
// error occurs, which is to call os.Exit(code). You can pass 'panic' as a function
// here if you prefer the panic() over os.Exit(1).
func BehaviorOnFatal(f func(string, int)) {
	fatalErrHandler = f
}

// fatal prints the message (if provided) and then exits.
func fatal(msg string, code int) {
	if msg != "" {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		_, _ = fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

// ErrExit may be passed to CheckError to instruct it to output nothing but exit with
// status code 1.
var ErrExit = fmt.Errorf("exit")

// CheckErr prints a user-friendly error to STDERR and exits with a non-zero
// exit code. Unrecognized errors will be printed with an "error: " prefix.
func CheckErr(err error) {
	checkErr(err, fatalErrHandler)
}

// checkErr formats a given error as a string and calls the passed handleErr
// func with that string and an exit code.
func checkErr(err error, handleErr func(string, int)) {
	if err == nil {
		return
	}

	switch {
	case err == ErrExit:
		handleErr("", DefaultErrorExitCode)
	//case clientcmd.IsConfigurationInvalid(err):
	//	handleErr(MultilineError("Error in configuration: ", err), DefaultErrorExitCode)
	default:
		msg := err.Error()
		if !strings.HasPrefix(msg, "error: ") {
			msg = fmt.Sprintf("error: %s", msg)
		}
		handleErr(msg, DefaultErrorExitCode)
	}
}
