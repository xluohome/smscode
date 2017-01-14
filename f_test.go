package main

import (
	"testing"

	//"github.com/issue9/assert"
)

func TestYuntongxun(t *testing.T) {

	var sms = NewSms()
	sms.SetServiceConfig("getpwd")
	sms.Code = "888888"
	sms.Mobile = "13575566313"

	//var y = &Yuntongxun{sms: sms}

	//y.Send()
}
