package jsonrpc

import (
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Client struct {
	connection *rpc.Client
}

func MakeClient(address string) (*Client, error) {
	 connection, err := jsonrpc.Dial("tcp", address)
	 if err != nil {
	 	return nil, err
	 }
	 return &Client{connection}, err
}

func (client *Client) Call(serviceName string, args interface{}, result *interface{}) error {
	defer client.connection.Close()
	return client.connection.Call(serviceName, args, result)
}