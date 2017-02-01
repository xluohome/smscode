package main

import (
	"testing"

	//"github.com/issue9/assert"
)

func TestAlidayu(t *testing.T) {

	var sms = NewSms()
	sms.SetServiceConfig("register")
	sms.Code = "888888"
	sms.Mobile = "13575566313"

	var y = &Alidayu{}

	if err := y.Send(sms); err != nil {
		t.Error(err)
	}

}

func TestYuntongxun(t *testing.T) {

	var sms = NewSms()
	sms.SetServiceConfig("getpwd")
	sms.Code = "888888"
	sms.Mobile = "13575566313"

	var y = &Yuntongxun{}

	if err := y.Send(sms); err != nil {
		t.Error(err)
	}

}

func TestWxhy(t *testing.T) {

	var sms = NewSms()
	sms.SetServiceConfig("getpwd")
	sms.Code = "888888"
	sms.Mobile = "13575566313"

	var y = &Hywx{}

	if err := y.Send(sms); err != nil {
		t.Error(err)
	}

}
