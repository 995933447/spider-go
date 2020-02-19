package url

import "strings"

func GetUriPath(url string) string {
	index := strings.LastIndex(url, "/")
	if index > 0 {
		return url[0:index + 1]
	}
	return url
}