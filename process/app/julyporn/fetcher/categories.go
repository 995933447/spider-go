package fetcher

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"util/htmlutil"
)

func Categories(url string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.AddCookie(&http.Cookie{
		Name: "l",
		Value: "zh",
	})

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
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

func DoFileWork()  {

}
