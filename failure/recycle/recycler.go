package recycler

import (
	"fmt"
	"spider-go/config"
	"spider-go/logger"
	"time"
)

var (
	FailureQueue []config.Request
	tick = time.Tick(time.Hour)
)

func RecycleFailedRequest() chan config.Request {
	failedChan := make(chan config.Request)
	go func() {
		for {
			failedRequest := <- failedChan
			FailureQueue = append(FailureQueue, failedRequest)
		}
	}()
	return failedChan
}

func Run() {
	go func() {
		<- tick
		for _, failedRequest := range FailureQueue {
			logger.DefaultLogger.Debug(fmt.Sprintf("Start do failed request:%s", failedRequest), nil)
		}
	}()
}