package app

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/logger"
)

func InitLogging() *logger.Logger {
	// Logging
	y, month, d := time.Now().Date()
	logName := strconv.Itoa(y) + "_" + month.String() + "_" + strconv.Itoa(d)
	f, err := os.OpenFile(logFile+logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) // 0666 | 0660 ?
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	l := logger.Init("Logger", true, true, f)
	logger.Info("********APP STARTED********")
	return l
}
