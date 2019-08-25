package fetcher

import (
	"spider-go/fetcher/config"
	"time"
)

var RateLimit = time.Tick(time.Second / config.Qps)