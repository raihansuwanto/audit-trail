package logger

import (
	"context"
	"io"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type defaultFields struct {
	fileAndLine interface{}
	fn          interface{}
	pkgPath     interface{}
}

var (
	defaultLogger = logrus.New()
)

func init() {

	defaultLogger.SetFormatter(&logrus.JSONFormatter{})
}

func SetOutput(output io.Writer) {
	defaultLogger.SetOutput(output)
}

func SetLogLevel(level logrus.Level) {
	defaultLogger.SetLevel(level)
}

func SetFormatter(formatter logrus.Formatter) {
	defaultLogger.SetFormatter(formatter)
}

func I() *logrus.Entry {
	val := getDefaultFieldValue()
	return defaultLogger.WithFields(logrus.Fields{
		"file":     val.fileAndLine,
		"function": val.fn,
		"pkgPath":  val.pkgPath,
	})
}

func IWithTraceId(ctx context.Context) *logrus.Entry {
	val := getDefaultFieldValue()
	return defaultLogger.WithFields(logrus.Fields{
		"file":     val.fileAndLine,
		"function": val.fn,
		"pkgPath":  val.pkgPath,
		"traceId":  ctx.Value("traceId"),
	})
}

func getDefaultFieldValue() defaultFields {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		logrus.Error("Could not get context info for logger!")
		return defaultFields{}
	}
	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	callerFunc := runtime.FuncForPC(pc)
	funcForPC := callerFunc.Name()
	fn := funcForPC[strings.LastIndex(funcForPC, ".")+1:]
	pkg := funcForPC[:strings.LastIndex(funcForPC, ".")]
	return defaultFields{
		fileAndLine: filename,
		fn:          fn,
		pkgPath:     pkg,
	}
}
