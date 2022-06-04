package utils

import (
	"github.com/yanzijie/ylog"
)

var Log *ylog.ConsoleLog

func init() {
	Log = ylog.NewConsoleLogger(ylog.DebugLevel)
}
