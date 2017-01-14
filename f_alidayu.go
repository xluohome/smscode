package main

import (
	"fmt"

	log "github.com/golang/glog"
	"github.com/opensource-conet/alidayu"
)

type Alidayu struct {
	sms *SMS
}

/**
阿里大鱼 业务流控触发 条件：

 短信验证码 ：使用同一个签名，对同一个手机号码发送短信验证码，支持1条/分钟，累计7条/小时   【手机号，签名】
 短信通知： 使用同一个签名和同一个短信模板ID，对同一个手机号码发送短信通知，支持50条/日 【手机号，签名，模板id】

**/
func (a *Alidayu) Send() error {
	alidayu.Appkey = config.Vendors["alidayu"]["appkey"]
	alidayu.AppSecret = config.Vendors["alidayu"]["appSecret"]
	if config.Vendors["alidayu"]["issendbox"] == "true" {
		alidayu.IsDebug = true
	} else {
		alidayu.IsDebug = false
	}
	res, _ := alidayu.SendOnce(a.sms.Mobile, a.sms.Config.Signname, a.sms.Config.Tpl, `{"code":"`+a.sms.Code+`"}`)
	if !res.Success {
		log.V(1).Infof("Alidayu:%+v", res.ResultError)
		return fmt.Errorf("%s", res.ResultError.SubMsg)
	}
	return nil
}
