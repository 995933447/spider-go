package htmlutil

import (
	"bufio"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
)

func DetermineEncoding(r *bufio.Reader) (encoding.Encoding, error) {
	bytes, err := r.Peek(1024)
	if err != nil{
		return nil, err
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e, nil
}