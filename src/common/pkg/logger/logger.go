package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
)

type Logger struct {
	NoLoc     bool
	LocOffset int
}

var (
	Default = &Logger{}
)

func (l *Logger) output(logType logType, message string) {
	var (
		codeLocation string // relative path starting with src
		funcName     string
	)
	if !(l.NoLoc) {
		// (ref.) [How do you get a Golang program to print the line number of the error it just called?](https://stackoverflow.com/a/24809646)
		if pc, file, line, ok := runtime.Caller(2 + l.LocOffset); ok { // 2 for skipping this function (`output`) and the API (`Print`/`Fatal`)
			// `file` should be `/go/src/...` in the docker container, but we want to show as `src/...`
			fileTokens := strings.SplitN(file, "/", 3)
			codeLocation = fmt.Sprintf("%s:%d", fileTokens[2], line)

			// (ref.) [How to get the current function name](https://stackoverflow.com/questions/25927660/how-to-get-the-current-function-name)
			// (ref.) [Return last item of strings.Split() slice in Golang](https://stackoverflow.com/questions/50311213/return-last-item-of-strings-split-slice-in-golang)
			frames := runtime.CallersFrames([]uintptr{pc})
			frame, _ := frames.Next()
			funcName = frame.Function[strings.LastIndex(frame.Function, "/")+1:]
		} else {
			codeLocation = "???"
			funcName = "???"
		}
	}

	if b, err := json.Marshal(&struct {
		LogType      string `json:"type"`
		CodeLocation string `json:"loc,omitempty"`
		FuncName     string `json:"func,omitempty"`
		Message      string `json:"message"`
	}{
		LogType:      string(logType),
		CodeLocation: codeLocation,
		FuncName:     funcName,
		Message:      message,
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Fprintln(os.Stderr, string(b))
	}
}

func (l *Logger) Printf(format string, a ...any) {
	l.output(logTypeInfo, fmt.Sprintf(format, a...))
}

func (l *Logger) Fatal(err error) {
	l.output(logTypeError, err.Error())
	os.Exit(1)
}
