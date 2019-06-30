package common

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"mahjong/app/config"

	"github.com/google/logger"
)

// Logging
func InitLogging() *logger.Logger {
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

	l := logger.Init("Logger", true, true, f)
	logger.Info("********APP STARTED********")
	return l
}
