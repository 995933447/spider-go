package fetcher

import (
	"spider-go/fetcher"
	"time"
)

var rateLimit = time.Tick(time.Second/3000)

func Video(url string) ([]byte, error) {
	<- rateLimit
	return fetcher.DoHtmlWork(url, true)
}