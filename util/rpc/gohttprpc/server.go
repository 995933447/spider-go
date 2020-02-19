package gohttprpc

import (
	"net/http"
	"net/rpc"
)

func MakeServer(service interface{}, address string) error {
	err := rpc.Register(service)
	if err != nil {
		return err
	}
	rpc.HandleHTTP()
	return http.ListenAndServe(address, nil)
}