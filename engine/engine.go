package engine

import (
	"spider-go/config"
	recycler "spider-go/failure/recycle"
	"spider-go/logger"
)

type Engine struct {
	Schedule config.Schedule
	WorkerNum int
	FailedRequestRecycler chan config.Request
}

func (engine Engine) Run(seeds ...config.Request) {
	engine.Schedule.Run()

	for _, seed := range seeds {
		engine.Schedule.Submit(seed)
	}

	out := make(chan config.ParseResult)
	for i := 0; i < engine.WorkerNum; i++ {
		engine.createWorker(make(chan config.Request), out)
	}

	recycler.Run()

	for {
		parseResult := <- out
		for _, request := range parseResult.Requests {
			engine.Schedule.Submit(request)
		}
	}
}

func (engine Engine) createWorker(in chan config.Request, out chan config.ParseResult) {
	go func() {
		for {
			engine.Schedule.WorkerReady(in)
			request := <- in
			parseResult, err := Work(request)
			if err != nil {
				engine.FailedRequestRecycler <- request
				logger.DefaultLogger.Error(err, nil)
			}
			out <- parseResult
		}
	}()
}