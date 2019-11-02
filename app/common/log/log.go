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

var (
	globalLogger Logger
	depth        = 1
)

type Logger struct {
	gl  *logger.Logger
	std bool
}

func Info(v ...interface{}) {
	if globalLogger.std {
		log.Print(v...)
		return
	}
	globalLogger.gl.InfoDepth(depth, v...)
}

func Warning(v ...interface{}) {
	if globalLogger.std {
		log.Print(v...)
		return
	}
	globalLogger.gl.WarningDepth(depth, v...)
}

func Error(v ...interface{}) {
	if globalLogger.std {
		log.Print(v...)
		return
	}
	globalLogger.gl.ErrorDepth(depth, v...)
}

func Fatal(v ...interface{}) {
	if globalLogger.std {
		log.Print(v...)
		return
	}
	globalLogger.gl.FatalDepth(depth, v...)
}

func Infof(format string, v ...interface{}) {
	if globalLogger.std {
		log.Printf(format, v...)
		return
	}
	globalLogger.gl.InfoDepth(depth, fmt.Sprintf(format, v...))
}

func Warningf(format string, v ...interface{}) {
	if globalLogger.std {
		log.Printf(format, fmt.Sprintf(format, v...))
		return
	}
	globalLogger.gl.WarningDepth(depth, v...)
}

func Errorf(format string, v ...interface{}) {
	if globalLogger.std {
		log.Printf(format, fmt.Sprintf(format, v...))
		return
	}
	globalLogger.gl.ErrorDepth(depth, v...)
}

func Fatalf(format string, v ...interface{}) {
	if globalLogger.std {
		log.Printf(format, v...)
		return
	}
	globalLogger.gl.FatalDepth(depth, fmt.Sprintf(format, v...))
}

// Logging
func InitLogging() Logger {
	if os.Getenv("DOCKER_RUN") == "true" {
		globalLogger = Logger{std: true}
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

func (l Logger) Close() {
	if !l.std {
		globalLogger.gl.Close()
	}
}
