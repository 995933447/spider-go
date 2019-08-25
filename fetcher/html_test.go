package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestHtml(t *testing.T) {
	//resp, err := http.Get("http://www.baidu.com")
	//c, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Println(string(c))


	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://115.85.40.130")
	}
	transport := &http.Transport{Proxy: proxy}
	cli := &http.Client{Transport: transport,Timeout:30*time.Second}
	resp, err := HttpGet(cli, "https://www.xvideos.com/video49918481/japanese_eager_beaver")
	//resp, err := http.Get("https://www.xvideos.com/video49918481/japanese_eager_beaver")
	if err != nil {
		t.Error(err)
		return
	}
	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(c))


}