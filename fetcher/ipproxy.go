package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"spider-go/fetcher/config"
)

type IpProxy struct {
	Ip string
	Port int
	Location string
	Source string
	Speed int
}

func GetProxyIp() (string, error) {
	var ipProxy IpProxy
	resp, err := http.Get("http://" + config.IpProxyHost + "/get")
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(data, &ipProxy); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", ipProxy.Ip, ipProxy.Port), nil
}