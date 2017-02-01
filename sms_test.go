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

func TestSendSmsok(t *testing.T) { //浙江省内号码

	t.Log("【浙江省内号码  允许发送】")
	sms := NewSms()
	sms.SetServiceConfig("register")
	if err := sms.Send("13575566313"); err != nil {

		t.Fatal(err)
	}

	if err := sms.CheckCode(sms.Mobile, sms.Code); err != nil {

		t.Fatal(err)
	}

	t.Log(sms.Mobile, sms.Code)

}

func TestSendSmsfaild(t *testing.T) { //浙江省外的号码

	t.Log("【浙江省外的号码 不允许发送】")
	sms := NewSms()
	sms.SetServiceConfig("register")
	if err := sms.Send("15970772900"); err != nil {

		t.Fatal(err)
	}

	if err := sms.CheckCode(sms.Mobile, sms.Code); err != nil {

		t.Fatal(err)
	}

	t.Log(sms.Mobile, sms.Code)
}

func TestSendSmsfaild2(t *testing.T) { //不同的验证码服务, 因为签名相同 会触发 大鱼流控

	t.Log("【不同的验证码服务, 因为签名相同会触发大鱼流控】")
	sms := NewSms()
	sms.SetServiceConfig("restpwd")
	if err := sms.Send("13575566313"); err != nil {

		t.Fatal(err)
	}

	if err := sms.CheckCode(sms.Mobile, sms.Code); err != nil {

		t.Fatal(err)
	}

	t.Log(sms.Mobile, sms.Code)

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
