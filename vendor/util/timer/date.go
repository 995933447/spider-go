package timer

import (
	"strings"
	"time"
)

func NowDate(format string) string {
	replaceMap := map[string]string{
		"Y": "2006",
		"m": "01",
		"d": "02",
		"H": "03",
		"i": "04",
		"s": "05",
	}
	for key, replace := range replaceMap {
		format = strings.ReplaceAll(format, key, replace)
	}
	return time.Unix(time.Now().Unix(), -1).Format(format)
}
