package fetcher

import (
	"net/http"
	"net/url"
	"time"
)

func NewHttpClient() (*http.Client, error) {
	proxyIp, err := GetProxyIp()
	if err != nil {
		return nil, err
	}
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://" + proxyIp)
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport,Timeout:30*time.Second}
	return client, nil
}

func HttpGet(client *http.Client, url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(request)
}