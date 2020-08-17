package main

import (
	"time"

	"github.com/EricChiou/logger"
)

func main() {
	logger.Init("log/")

	logger.Trace.Println("some thing failed")
	logger.Info.Println("some thing failed")
	logger.Warn.Println("some thing failed")
	logger.Error.Println("some thing failed")

	time.Sleep(60 * time.Second)
}
