package main

import (
	"testing"

	"github.com/issue9/assert"
)

func TestConvert(t *testing.T) {

	a := assert.New(t)

	var x int64 = 1483680227

	b := Int64ToBytes(x)

	var i int64 = BytesToInt64(b)

	a.Equal(x, i, true)
}
