package surl

import (
	"fmt"
	"testing"
)

func TestSurl(t *testing.T) {
	c := NewCreator(LenSmall)

	c.SetKey(c.NewKey())

	for i := 1; i < 10; i++ {
		s, _ := c.Encode(int64(i))
		n, _ := c.Decode(s)
		fmt.Println(i, "\t", s, "\t", n)
	}
}
