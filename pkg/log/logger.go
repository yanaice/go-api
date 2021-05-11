package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/google/uuid"
	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
	"go-starter-project/pkg/threadlocal"
)

type Logger struct {
	logger     *logrus.Logger
	env        string
	service    string
	instanceId string
}

type Fields map[string]interface{}

var (
	std = &Logger{
		logger:     new(),
		env:        "startup",
		service:    "startup",
		instanceId: uuid.New().String(),
	}
)

const (
	minimumCallerDepth = 2
	maximumCallerDepth = 15
)

var excludePkgName = map[string]bool{
	"gitlab.com/diancai/diancai-services-common/log":    true,
	"gitlab.com/diancai/diancai-services-common/derror": true,
}

func new() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(joonix.NewFormatter())
	return logger
}

func parseLevel(lvl string) (logrus.Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return logrus.PanicLevel, nil
	case "fatal":
		return logrus.FatalLevel, nil
	case "error":
		return logrus.ErrorLevel, nil
	case "warn", "warning":
		return logrus.WarnLevel, nil
	case "info":
		return logrus.InfoLevel, nil
	case "debug":
		return logrus.DebugLevel, nil
	case "trace":
		return logrus.TraceLevel, nil
	}

	var l logrus.Level
	return l, fmt.Errorf("not a valid logrus Level: %q", lvl)
}

func wrapFields() *logrus.Entry {
	file, line := GetCaller()

	id := threadlocal.GetCorrelationID()

	return std.logger.WithFields(logrus.Fields{
		"instance-id":    std.instanceId,
		"env":            std.env,
		"service":        std.service,
		"file":           file,
		"line":           line,
		"correlation-id": id,
	})
}

func GetCaller() (string, int) {
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if !excludePkgName[pkg] {
			return f.Function, f.Line
		}
	}

	return "", 0
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

func Init(env, service string, level string) {
	lv, err := parseLevel(level)
	if err != nil {
		log.Fatal(err)
	}

	std.logger.SetLevel(lv)
	std.logger.SetFormatter(joonix.NewFormatter())

	std.env = env
	std.service = service
}

func InitWithFile(env, service string, level string, filename string) {
	Init(env, service, level)

	logFile, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	std.logger.SetOutput(logFile)
}

func Logrus() *logrus.Logger {
	return std.logger
}

func GetLevel() string {
	return std.logger.GetLevel().String()
}

func IsLevelEnabled(level string) (bool, error) {
	lv, err := parseLevel(level)
	if err != nil {
		return false, err
	}

	return std.logger.IsLevelEnabled(lv), nil
}

// WithContext creates an entry from the standard logger and adds a context to it.
func WithContext(ctx context.Context) *logrus.Entry {
	return wrapFields().WithContext(ctx)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *logrus.Entry {
	return wrapFields().WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields Fields) *logrus.Entry {
	data := make(logrus.Fields, len(fields))
	for k, v := range fields {
		data[k] = v
	}

	return wrapFields().WithFields(data)
}

func Trace(args ...interface{}) {
	wrapFields().Trace(args...)
}

func Debug(args ...interface{}) {
	wrapFields().Debug(args...)
}

func Print(args ...interface{}) {
	wrapFields().Print(args...)
}

func Info(args ...interface{}) {
	wrapFields().Info(args...)
}

func Warn(args ...interface{}) {
	wrapFields().Warn(args...)
}

func Warning(args ...interface{}) {
	wrapFields().Warning(args...)
}

func Error(args ...interface{}) {
	wrapFields().Error(args...)
}

func Panic(args ...interface{}) {
	wrapFields().Panic(args...)
}

func Fatal(args ...interface{}) {
	wrapFields().Fatal(args...)
}

func Tracef(format string, args ...interface{}) {
	wrapFields().Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	wrapFields().Debugf(format, args...)
}

func Printf(format string, args ...interface{}) {
	wrapFields().Printf(format, args...)
}

func Infof(format string, args ...interface{}) {
	wrapFields().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	wrapFields().Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	wrapFields().Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	wrapFields().Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	wrapFields().Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	wrapFields().Fatalf(format, args...)
}

func Traceln(args ...interface{}) {
	wrapFields().Traceln(args...)
}

func Debugln(args ...interface{}) {
	wrapFields().Debugln(args...)
}

func Println(args ...interface{}) {
	wrapFields().Println(args...)
}

func Infoln(args ...interface{}) {
	wrapFields().Infoln(args...)
}

func Warnln(args ...interface{}) {
	wrapFields().Warnln(args...)
}

func Warningln(args ...interface{}) {
	wrapFields().Warningln(args...)
}

func Errorln(args ...interface{}) {
	wrapFields().Errorln(args...)
}

func Panicln(args ...interface{}) {
	wrapFields().Panicln(args...)
}

func Fatalln(args ...interface{}) {
	wrapFields().Fatalln(args...)
}
