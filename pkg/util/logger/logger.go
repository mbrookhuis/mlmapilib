package logger

import (
	"fmt"
	"io"
	"os"

	model "github.com/mbrookhuis/mlmapilib/pkg/models/inputfile"

	"github.com/sirupsen/logrus"
	"strings"
	"sync"
)

// Logger is the application's central logger.
var Logger *logrus.Logger
var once sync.Once

// InitLogger initializes the global logger based on the application configuration.
func InitLogger(genConfig model.ConfigGeneral) error {
	once.Do(func() {
		Logger = logrus.New()
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		// Configure screen output
		screenLevel, err := logrus.ParseLevel(genConfig.General.Log.ScreenLevel)
		if err != nil {
			fmt.Printf("Invalid screen log level '%s', defaulting to info: %v\n", genConfig.General.Log.ScreenLevel, err)
			screenLevel = logrus.InfoLevel
		}
		Logger.SetOutput(os.Stdout) // Default output to screen
		Logger.SetLevel(screenLevel)
		// Configure file output if a file path is provided
		if genConfig.General.Log.FilePath != "" {
			fileLevel, err := logrus.ParseLevel(genConfig.General.Log.FileLevel)
			if err != nil {
				fmt.Printf("Invalid file log level '%s', defaulting to debug: %v\n", genConfig.General.Log.FileLevel, err)
				fileLevel = logrus.DebugLevel
			}
			logDir := ""
			lastSlash := strings.LastIndex(genConfig.General.Log.FilePath, "/")
			if lastSlash != -1 {
				logDir = genConfig.General.Log.FilePath[:lastSlash]
				if err := os.MkdirAll(logDir, 0755); err != nil {
					fmt.Printf("Failed to create log directory '%s': %v\n", logDir, err)
					return
				}
			}
			logFile, err := os.OpenFile(genConfig.General.Log.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				fmt.Printf("Failed to open log file '%s': %v\n", genConfig.General.Log.FilePath, err)
				return
			}
			mw := io.MultiWriter(os.Stdout, logFile)
			Logger.SetOutput(mw)
			if fileLevel < screenLevel {
				Logger.SetLevel(fileLevel)
			} else {
				Logger.SetLevel(screenLevel)
			}
		}
	})
	return nil
}

func Debug(args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.DebugLevel) {
		Logger.Debug(args...)
	}
}

func Info(args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.InfoLevel) {
		Logger.Info(args...)
	}
}

func Warn(args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.WarnLevel) {
		Logger.Warn(args...)
	}
}

func Error(args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.ErrorLevel) {
		Logger.Error(args...)
	}
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.DebugLevel) {
		Logger.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.InfoLevel) {
		Logger.Infof(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.WarnLevel) {
		Logger.Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.ErrorLevel) {
		Logger.Errorf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// NOTE: For truly separate log levels for screen and file,
// you would typically need to implement `logrus.Hook` for each output.
// A simpler approach for the scope of this example is to set the
// *global* log level of the Logrus instance to the most permissive
// of the two (e.g., if screen is INFO and file is DEBUG, the global
// level should be DEBUG). Then, within the hook or before writing
// to the specific output, you would filter based on the desired level
// for that output.
