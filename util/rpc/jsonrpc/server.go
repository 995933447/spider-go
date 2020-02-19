package jsonrpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	writer "util/logger/writters/console"
)

func MakeServer(service interface{}, address string) error {
	err := rpc.Register(service)
	if err != nil {
		return err
	}

	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	for {
		connection, err := listen.Accept()
		if err != nil {
			err := writer.ConsoleWriter{}.Error(err, nil)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		go jsonrpc.ServeConn(connection)
	}
}