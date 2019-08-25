package logger

import (
	"util/logger"
	cwriter "util/logger/writers/console"
	fwriter "util/logger/writers/file"
)

var (
	DefaultLogger = logger.Logger{}
	fw = fwriter.FileWriter{}
	cw = cwriter.ConsoleWriter{}
)

func init() {
	DefaultLogger.RegisterWriter(fw)
	DefaultLogger.RegisterWriter(cw)
}