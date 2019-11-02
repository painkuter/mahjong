package log

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/logger"

	"mahjong/app/config"
)

var globalLogger Logger

type Logger struct {
	gl  *logger.Logger
	std bool
}

func Info(v ...interface{}) {
	if globalLogger.std {
		log.Print(v...)
	}
	globalLogger.gl.Info(v...)
}

func Warning(v ...interface{}) {
	if globalLogger.std {
		log.Print(v...)
	}
	globalLogger.gl.Warning(v...)
}

func Error(v ...interface{}) {
	if globalLogger.std {
		log.Print(v...)
	}
	globalLogger.gl.Error(v...)
}

func Fatal(v ...interface{}) {
	if globalLogger.std {
		log.Print(v...)
	}
	globalLogger.gl.Fatal(v...)
}

func Infof(format string, v ...interface{}) {
	if globalLogger.std {
		log.Printf(format, v...)
	}
	globalLogger.gl.Info(v...)
}

func Warningf(format string, v ...interface{}) {
	if globalLogger.std {
		log.Printf(format, v...)
	}
	globalLogger.gl.Warning(v...)
}

func Errorf(format string, v ...interface{}) {
	if globalLogger.std {
		log.Printf(format, v...)
	}
	globalLogger.gl.Error(v...)
}

func Fatalf(format string, v ...interface{}) {
	if globalLogger.std {
		log.Printf(format, v...)
	}
	globalLogger.gl.Fatal(v...)
}

// Logging
func InitLogging() Logger {
	if os.Getenv("DOCKER_RUN") == "true" {
		return globalLogger
	}

	// checking directory exists
	if _, err := os.Stat(config.LogDir); os.IsNotExist(err) {
		os.MkdirAll(config.LogDir, os.ModePerm)
	}

	y, month, d := time.Now().Date()
	logName := strconv.Itoa(y) + "_" + month.String() + "_" + strconv.Itoa(d)
	f, err := os.OpenFile(config.LogDir+"/"+config.LogPrefix+logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) // 0666 | 0660 ?
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	globalLogger.gl = logger.Init("Logger", true, false, f)
	log.Printf("********APP STARTED********")
	return globalLogger
}

func (l *Logger) Close() {
	if !l.std {
		globalLogger.gl.Close()
	}
}
