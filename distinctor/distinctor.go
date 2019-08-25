package distinctor

import (
	"fmt"
	"spider-go/persist/redis"
	"time"
	"util/timer"
)

type Distinctor struct {
	recordChan chan string
}

func setFetchedRecorder(distinctor *Distinctor) {
	distinctor.recordChan = make(chan string)
}

func Run(distinctor *Distinctor) error {
	setFetchedRecorder(distinctor)

	redisClient, err := redis.NewClient()
	if err != nil {
		return err
	}

	go func() {
		for {
			url := <- distinctor.recordChan
			redisClient.SetNX(url, timer.NowDate("Y-m-d H:i:s"), 0)
		}
	}()
	return nil
}

func (distinctor *Distinctor) RecordFetched(url string) {
	distinctor.recordChan <- url
}

func (*Distinctor) CheckIsFetched(url string) (bool, error) {
	redisClient, err := redis.NewClient()
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	result, err := redisClient.Get(url).Result()
	if err != nil {
		return false, nil
	}

	if _, err := time.Parse("2006-01-02 15:04:05", result); err != nil {
		return false, nil
	}

	return true, nil
}