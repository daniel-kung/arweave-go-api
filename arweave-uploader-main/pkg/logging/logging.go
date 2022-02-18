package logging

import (
	"os"
	"time"

	"ccian.cc/really/arweave-api/pkg/setting"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	fileMode = "file"
	stdMode  = "stdout"
)

var (
	gLogger   *logrus.Logger
	gLogEntry *logrus.Entry
)

func Setup(config *setting.LoggerConfig) {
	gLogger = logrus.New()
	gLogger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	if config.EnableContext {
		gLogger.SetReportCaller(true)
	}

	if level, err := logrus.ParseLevel(config.Level); err != nil {
		gLogger.SetLevel(logrus.InfoLevel)
	} else {
		gLogger.SetLevel(level)
	}

	switch config.Mode {
	case fileMode:
		path := config.File.NamePrefix
		writer, err := rotatelogs.New(
			path+"-%Y%m%d.log",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
			rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
		)
		if err != nil {
			logrus.Fatalf("logging file open failed: %v", err)
		}

		hook := lfshook.NewHook(
			addWriterToLowestLevel(writer, gLogger.GetLevel()),
			&logrus.JSONFormatter{},
		)
		logrus.StandardLogger().AddHook(hook) // logrus global logger
		gLogger.AddHook(hook)                 // our global logger
	case stdMode:
		fallthrough
	default:
		gLogger.SetOutput(os.Stdout)
	}

	gLogEntry = NewLogEntry("global")
}

// default lowest level is info
func addWriterToLowestLevel(writer *rotatelogs.RotateLogs, lowestLevel logrus.Level) lfshook.WriterMap {
	switch lowestLevel {
	case logrus.DebugLevel:
		return lfshook.WriterMap{
			logrus.DebugLevel: writer,
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		}
	case logrus.WarnLevel:
		return lfshook.WriterMap{
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		}
	case logrus.ErrorLevel:
		return lfshook.WriterMap{
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		}
	case logrus.FatalLevel:
		return lfshook.WriterMap{
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		}
	case logrus.PanicLevel:
		return lfshook.WriterMap{
			logrus.PanicLevel: writer,
		}
	case logrus.InfoLevel:
		fallthrough
	default:
		return lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		}
	}
}

func NewLogEntry(moduleName string) *logrus.Entry {
	return gLogger.WithField("module", moduleName)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return gLogEntry.WithFields(fields)
}
