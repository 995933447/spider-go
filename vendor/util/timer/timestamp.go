package timer

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var DefaultTimeZone = time.FixedZone("CST", 8 * 3600)

func DurationToSecond(duration string) (int, error) {
	durationSegments := strings.Split(duration, ":")

	var hours, minutes, seconds int
	switch len(durationSegments) {
	case 3:
		hours, _ = strconv.Atoi(durationSegments[0])
		minutes, _ = strconv.Atoi(durationSegments[1])
		seconds, _ = strconv.Atoi(durationSegments[2])
	case 2:
		hours = 0
		minutes, _ = strconv.Atoi(durationSegments[0])
		seconds, _ = strconv.Atoi(durationSegments[1])
	default:
		return -1, errors.New("please input arg formatter like : 00:00:00")
	}

	return hours * 3600 + minutes * 60 + seconds, nil
}

func NowUnix(timezone *time.Location) int {
	if timezone == nil {
		return int(time.Now().In(DefaultTimeZone).Unix())
	}
	return int(time.Now().In(timezone).Unix())
}