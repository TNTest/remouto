package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestConcat(t *testing.T) {
	a := []byte{5, 5, 5}
	b := []byte{6, 6, 6}
	c := concat(a, b)
	d := []byte{5, 5, 5, 6, 6, 6}
	fmt.Print(c)
	if !bytes.Equal(c, d) {
		t.Fail()
	}
}
