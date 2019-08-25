package itemworker

import (
	"fmt"
	path2 "path"
	"strconv"
	"strings"
	"testing"
)

func TestVideos(t *testing.T) {
	path := "/a/b/abc.go"
	fmt.Println(path2.Dir(path))
	long := "3.1342"
	if n, err := strconv.ParseFloat(long, 64); err != nil {
		t.Error(err)
	} else {
		fmt.Println(n)
	}
	fmt.Println(strings.TrimRight(path2.Dir("654.ts"), "."))
}
