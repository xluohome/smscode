package main

import (
	"testing"
)

func BenchmarkSendSms111(b *testing.B) {

	for i := 0; i < b.N; i++ {

		sms := NewSms()
		sms.SetServiceConfig("register")
		if err := sms.Send("13575566313"); err != nil {
			b.Log(err)
			//b.Error(err)
		}

	}
}

var (
	test_mobile  = "13575566313"
	test_mobile2 = "13375566310"
	test_code    = ""
)

func TestCheckAreaOk(t *testing.T) {

	sms := NewSms()
	sms.SetServiceConfig("register")
	sms.Config.Allowcity = []string{"0575"}
	sms.Mobile = "13575566310" //绍兴的手机号码
	if err := sms.checkArea(); err != nil {
		t.Error("归属地限制无效")
	}
}

func TestCheckAreaInvild(t *testing.T) {

	sms := NewSms()
	sms.SetServiceConfig("register")
	sms.Config.Allowcity = []string{"0575"}
	sms.Mobile = "13375566310" //不是绍兴的手机号码
	if err := sms.checkArea(); err != nil {
		t.Log(err)
		return
	}
	t.Error("归属地限制无效")
}

func TestSendSmsok(t *testing.T) {
	sms := NewSms()
	sms.SetServiceConfig("register")
	sms.Config.Allowcity = []string{"0575"}
	if err := sms.Send(test_mobile); err != nil {
		t.Error(err)
		return
	}
	test_code = sms.Code
	t.Run("CheckCode", testCheckCode)
}

func testCheckCode(t *testing.T) {
	sms := NewSms()
	sms.SetServiceConfig("register")
	if err := sms.CheckCode(test_mobile, test_code); err != nil {
		t.Error(err)
	}
}

func TestSendSmsfaild(t *testing.T) { //绍兴外的号码

	sms := NewSms()
	sms.SetServiceConfig("register")
	sms.Config.Allowcity = []string{"0575"}
	if err := sms.Send(test_mobile2); err != nil {
		t.Log(err)
		return
	}
	t.Error("绍兴外的号码也能发送")
	return
}

func TestSendSmsinfo(t *testing.T) { //手机验证码发送信息
	t.Log("【显示手机验证码发送信息】")
	sms := NewSms()
	sms.SetServiceConfig("register")

	var infos, _ = sms.Info("13575566313")

	for k, v := range infos {

		t.Logf("%s:%v", k, v)
	}

}
