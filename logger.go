package main

import (
	"context"
	"fmt"
	"runtime"
)

const LOG_FILE = "./numbers.log"

var terminate = fmt.Sprintf("%s%s", "terminate", LINE_BREAK)
var currOs = runtime.GOOS

type Appender interface {
	AddToLog(ctx context.Context, path, input string) error
}
