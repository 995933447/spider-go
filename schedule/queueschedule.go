package schedule

import "spider-go/config"

type QueueSchedule struct {
	requestChan chan config.Request
	workerChan chan chan config.Request
}

func (schedule *QueueSchedule) Submit(request config.Request)  {
	go func() {
		schedule.requestChan <- request
	}()
}

func (schedule *QueueSchedule) WorkerReady(readyWorker chan config.Request) {
	schedule.workerChan <- readyWorker
}

func (schedule *QueueSchedule) Run() {
	schedule.requestChan = make(chan config.Request)
	schedule.workerChan = make(chan chan config.Request)

	requestQueue := []config.Request{}
	workerQueue := []chan config.Request{}

	go func() {
		for {
			var activeRequest config.Request
			var activeWorker chan config.Request
			if len(requestQueue) > 0 && len(workerQueue) > 0 {
				activeRequest = requestQueue[0]
				activeWorker = workerQueue[0]
			}
			select {
				case request := <- schedule.requestChan:
					requestQueue = append(requestQueue, request)
				case worker := <- schedule.workerChan:
					workerQueue = append(workerQueue, worker)
				case activeWorker <- activeRequest:
					requestQueue = requestQueue[1:]
					workerQueue = workerQueue[1:]
			}
		}
	}()
}