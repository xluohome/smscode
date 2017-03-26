package main

import (
	"testing"

	//"github.com/issue9/assert"
)

func TestGetConfig(t *testing.T) {

	var conf Config

	err := conf.ParseConfigFile("conf/conf.yaml")

	if err != nil {

		t.Error(err)
	}

	t.Log(conf.Vendors["yuntongxun"])

}
