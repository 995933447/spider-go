package url

import (
	gourl "net/url"
)

func GetHost(url string) (string, error) {
	info, err := gourl.Parse(url)
	if err != nil {
		return "", err
	}
	return info.Host, err
}

func GetScheme(url string) (string, error) {
	info, err := gourl.Parse(url)
	if err != nil {
		return "", err
	}

	return info.Scheme, nil
}