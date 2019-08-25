package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"spider-go/fetcher"
	"spider-go/fetcher/config"
	"spider-go/logger"
	"time"
	"util/rand"
	urltool "util/url"
)


func File(url string) ([]byte, error) {
	<- fetcher.RateLimit
	return DoFileWork(url, false)
}

func DoFileWork(url string, useProxy bool) (content []byte, err error) {
	logger.DefaultLogger.Info("Fetching url :" + url, nil)

	var resp *http.Response

	try := 0
	for try <= 10{
		httpClient := &http.Client{}
		var request *http.Request
		request, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}


		request.Header.Add("User-Agent", config.UserAgent[rand.RandIntN(len(config.UserAgent) - 1)])
		//request.Header.Add("X-Requested-With", "ShockwaveFlash/27.0.0.170")
		if host, err := urltool.GetHost(url); err != nil {
			request.Header.Add("Host", host)
		}
		resp, err = httpClient.Do(request)
		if err != nil {
			try++
			time.Sleep(time.Second * time.Duration(rand.RandIntN(20)))
				time.Sleep(time.Second * time.Duration(rand.RandIntN(10)))
			} else {
				break
			}
	}


	if err != nil {
		return nil ,fmt.Errorf("fetch url error:%s",err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch url %s occured error:%d", url, resp.StatusCode)
	}

	all ,err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("fetch url %s occured error:%s", url, err)
	}

	return all, nil
}