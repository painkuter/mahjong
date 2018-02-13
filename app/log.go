package app

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/logger"
)

// Logging
func InitLogging() *logger.Logger {
	// checking directory exists
	if _, err := os.Stat(logDir); os.IsNotExist(err){
		os.MkdirAll(logDir, os.ModePerm)
	}

	y, month, d := time.Now().Date()
	logName := strconv.Itoa(y) + "_" + month.String() + "_" + strconv.Itoa(d)
	f, err := os.OpenFile(logDir+"/"+logPrefix+logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666) // 0666 | 0660 ?
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	l := logger.Init("Logger", true, true, f)
	logger.Info("********APP STARTED********")
	return l
}
