package main

import (
	"testing"
	"time"

	"github.com/issue9/assert"
)

var (
	test2_mobile  = "13575566311"
	test2_smscode = "888888"
)

func TestMobileArea(t *testing.T) {

	a := assert.New(t)

	sms := NewSms()
	sms.SetServiceConfig("register")
	sms.Config.Allowcity = []string{"0575"}
	sms.Mobile = test2_mobile

	model := NewModel(sms)
	area, err := model.GetMobileArea()

	if err != nil {

		t.Error(err)
	}

	a.Equal(area, "0575")

}

func TestSendTime(t *testing.T) {

	a := assert.New(t)

	sms := NewSms()

	sms.SetServiceConfig("register")

	sms.Mobile = test2_mobile
	sms.Config.Group = "testG"
	sms.Config.Signname = "luoluo"

	model := NewModel(sms)

	model.SetSendTime()

	time1, err := model.GetSendTime()

	if err != nil {

		t.Error(err)
	}

	a.Equal(time1, sms.NowTime.Unix())

}

func TestTodaySendNums(t *testing.T) {

	a := assert.New(t)

	sms := NewSms()

	sms.SetServiceConfig("register")
	sms.Mobile = test2_mobile
	sms.Config.Group = "testG"
	sms.Config.Signname = "luoluo"

	model := NewModel(sms)

	model.SetTodaySendNums(1)

	num, err := model.GetTodaySendNums()
	if err != nil {

		t.Error(err)
	}

	a.Equal(num, 1)

}

func TestSmsCode(t *testing.T) {

	a := assert.New(t)

	sms := NewSms()

	sms.SetServiceConfig("register")
	sms.Config.Group = "testG"
	sms.Config.Signname = "luoluo"

	sms.Mobile = test2_mobile
	sms.Code = test2_smscode

	model := NewModel(sms)
	model.SetSmsCode()

	code, uxtime, err := model.GetSmsCode()
	if err != nil {
		t.Error(err)
	}

	if time.Now().Unix() > uxtime {
		a.False(false)
	}

	a.Equal(code, test2_smscode)
}
