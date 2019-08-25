package fetcher

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"spider-go/logger"
	"util/htmlutil"
)

func Html(url string) (content []byte, err error) {
	<- RateLimit
	return DoHtmlWork(url, false)
}

func HtmlByProxy(url string) (content []byte, err error) {
	<-RateLimit
	return DoHtmlWork(url, true)
}

func DoHtmlWork(url string, useProxy bool) (content []byte, err error) {
	logger.DefaultLogger.Info("Fetching " + url, nil)

	var resp *http.Response
	if useProxy {
		var httpClient *http.Client
		httpClient, err = NewHttpClient()
		if err != nil {
			logger.DefaultLogger.Warning(err, nil)
			resp, err = http.Get(url)
		} else {
			resp, err = HttpGet(httpClient, url)
		}
	} else {
		resp, err = http.Get(url)
	}

	if err != nil {
		return nil, fmt.Errorf("fetch url %s occured error:%s", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch url %s occured error:%d", url, resp.StatusCode)
	}

	bufReader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()
	e, err := htmlutil.DetermineEncoding(bufReader)
	if err != nil {
		e = unicode.UTF8
	}

	utf8Reader := transform.NewReader(bufReader, e.NewDecoder())

	return ioutil.ReadAll(utf8Reader)
}