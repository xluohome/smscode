package main

import (
	"testing"
	"time"

	"github.com/issue9/assert"
)

func TestMobileArea(t *testing.T) {

	a := assert.New(t)

	sms := NewSms()
	sms.SetServiceConfig("register")

	sms.Mobile = "13575566313"

	model := NewModel(sms)

	model.SetMobileArea("0575")

	t.Log("13575566313 SetMobileArea  0575 ")

	area, err := model.GetMobileArea()

	if err != nil {

		t.Error(err)
	}

	a.Equal(area, "0575")

	t.Log("13575566313 GetMobileArea:", area)

}

func TestSendTime(t *testing.T) {

	a := assert.New(t)

	sms := NewSms()

	sms.SetServiceConfig("register")

	sms.Mobile = "13575566313"

	model := NewModel(sms)

	model.SetSendTime()

	t.Log("13575566313 SetSendTime  ", sms.NowTime.Unix())

	time1, err := model.GetSendTime()

	if err != nil {

		t.Error(err)
	}

	a.Equal(time1, sms.NowTime.Unix())

	t.Log("13575566313 GetSendTime:", time1)

}

func TestTodaySendNums(t *testing.T) {

	a := assert.New(t)

	sms := NewSms()

	sms.SetServiceConfig("register")

	sms.Mobile = "13575566313"

	model := NewModel(sms)

	model.SetTodaySendNums(1)

	t.Log("13575566313 SetTodaySendNums  1 ")

	num, err := model.GetTodaySendNums()
	if err != nil {

		t.Error(err)
	}

	a.Equal(num, 1)

	t.Log("13575566313 GetTodaySendNums:", num)

}

func TestSmsCode(t *testing.T) {

	a := assert.New(t)

	sms := NewSms()

	sms.SetServiceConfig("register")

	sms.Mobile = "13575566313"

	sms.Code = "888888"

	model := NewModel(sms)

	model.SetSmsCode()

	t.Log("13575566313 SetSmsCode  888888 ", sms.NowTime.Unix())

	code, uxtime, err := model.GetSmsCode()
	if err != nil {

		t.Error(err)
	}

	if time.Now().Unix() > uxtime {
		a.False(false)
	}

	a.Equal(code, "888888")

	t.Log("13575566313 GetSmsCode:", code, uxtime)

}
