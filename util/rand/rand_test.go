package rand

import (
	"fmt"
	"testing"
)

func TestRand(t *testing.T)  {
	fmt.Println(RandInt())
	fmt.Println(RandIntN(500))
}
