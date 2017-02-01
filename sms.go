package main

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/golang/glog"
)

type Sender interface {
	Send(sms *SMS) error
}

func sendcode(sms *SMS) error {
	var vname = sms.Config.Vendor
	v, ok := smsvendor[vname]
	if !ok {
		log.Fatal("您设置的短信通道服务商有误")
		return fmt.Errorf("%s", "短信通道设置不正确")
	}
	return v.Send(sms)
}

var (
	smap      = make(map[string]bool)
	smsmut    sync.Mutex
	smsvendor = make(map[string]Sender)
)

type SMS struct {
	Mobile      string
	Code        string
	Uid         string
	serviceName string
	Config      *ServiceConfig
	ConfigisOK  bool
	NowTime     time.Time
	model       *Model
}

func NewSms() *SMS {
	sms := &SMS{}
	sms.model = NewModel(sms)
	sms.NowTime = time.Now()
	return sms
}

//设置服务配置文件
func (sms *SMS) SetServiceConfig(serviceName string) *SMS {
	smsmut.Lock()
	defer smsmut.Unlock()
	sms.Config, sms.ConfigisOK = config.ServiceList[serviceName]
	if sms.ConfigisOK {
		sms.serviceName = serviceName
	}
	return sms
}

//  归属地规则校验
func (sms *SMS) checkArea() error {

	if len(sms.Config.Allowcity) < 1 { //没有启用
		return nil
	}

	area, err := sms.model.GetMobileArea()
	if err != nil {
		return err
	}

	var Allow = false
	for _, citycode := range sms.Config.Allowcity {
		if strings.Contains(area, citycode) {
			Allow = true //允许发送sms
			break
		}
	}

	if !Allow {
		return fmt.Errorf(config.Errormsg["err_allow_areacode"], strings.Join(sms.Config.Allowcity, ","))
	}

	return nil
}

func (sms *SMS) checkhold() error {

	sendTime, err := sms.model.GetSendTime()
	if err != nil {
		return err
	}

	if sendTime > 0 && sms.NowTime.Unix()-sendTime < Maxsendtime { //发送间隔不能小于60秒
		return fmt.Errorf(config.Errormsg["err_per_minute_send_num"])
	}

	sendMax, err := sms.model.GetTodaySendNums()
	if err != nil {
		return err
	}

	if sendMax > 0 && sendMax >= sms.Config.MaxSendNums {
		return fmt.Errorf(config.Errormsg["err_per_day_max_send_nums"], sms.Config.MaxSendNums)
	}

	return nil
}

/**
当前模式  1：只有手机号对应的uid存在时才能发送，2：只有uid不存在时才能发送，3：不管uid是否存在都发送
**/
func (sms *SMS) currModeok() error {

	uid, err := sms.model.GetSmsUid()
	if err != nil {
		return err
	}
	switch mode := sms.Config.Mode; mode {
	case 0x01:
		if uid != "" {
			return nil
		}
		return fmt.Errorf(config.Errormsg["err_model_not_ok1"], sms.Mobile)
	case 0x02:
		if uid == "" {
			return nil
		}
		return fmt.Errorf(config.Errormsg["err_model_not_ok2"], sms.Mobile)
	case 0x03:
		return nil
	}

	return fmt.Errorf("请正确配置config中的mode参数")
}

//保存数据
func (sms *SMS) save() {

	sms.model.SetSendTime()

	nums, _ := sms.model.GetTodaySendNums()

	newnums := atomic.AddUint64(&nums, 1) //原子操作+1

	sms.model.SetTodaySendNums(newnums)

	sms.model.SetSmsCode()
}

//发送短信
func (sms *SMS) Send(mobile string) error {
	if !sms.ConfigisOK {
		return fmt.Errorf("(%s)服务配置不存在", sms.serviceName)
	}

	/**
	限制一个手机号只允许在一次send成功（失败）后再次send
	为什么这么做？
	1：在高并发下保证一个手机号的send操作是同步的，后续规则校验可以依次进行；
	2：同时保证高并发下的send性能；
	**/
	smsmut.Lock()
	if _, ok := smap[mobile]; ok {
		smsmut.Unlock()
		return fmt.Errorf("服务器忙请稍后再试...")
	}
	smap[mobile] = true
	smsmut.Unlock()

	//send返回后取消限制
	defer func() {
		smsmut.Lock()
		delete(smap, mobile)
		smsmut.Unlock()
	}()

	sms.Mobile = mobile

	//生成验证码
	sms.Code = makeCode()

	if err := VailMobile(sms.Mobile); err != nil {
		return err
	}
	if err := sms.checkArea(); err != nil {
		return err
	}
	if err := sms.currModeok(); err != nil {
		return err
	}
	if err := sms.checkhold(); err != nil {
		return err
	}
	if err := sendcode(sms); err != nil {

		//发送失败 callback
		AddCallbackTask(sms, "Failed")
		return err
	}

	//保存记录
	sms.save()

	//发送成功 callback
	AddCallbackTask(sms, "Success")

	return nil
}

func (sms *SMS) CheckCode(mobile, code string) error {
	if !sms.ConfigisOK {
		return fmt.Errorf("(%s)服务配置不存在", sms.serviceName)
	}

	sms.Mobile = mobile
	sms.Code = code

	if err := VailMobile(sms.Mobile); err != nil {
		return err
	}

	if err := VailCode(sms.Code); err != nil {
		return err
	}

	oldcode, validtime, _ := sms.model.GetSmsCode()

	if sms.Code != oldcode {
		return fmt.Errorf(config.Errormsg["err_code_not_ok"], sms.Code)
	}

	if sms.NowTime.Unix() > validtime {
		time1 := time.Unix(validtime, 0)
		return fmt.Errorf(config.Errormsg["err_vailtime_not_ok"], time.Since(time1).String())

	}

	//验证成功时 callback
	AddCallbackTask(sms, "Checkok")

	return nil
}

func (sms *SMS) SetUid(mobile, uid string) error {
	if !sms.ConfigisOK {
		return fmt.Errorf("(%s)服务配置不存在", sms.serviceName)
	}

	sms.Mobile = mobile
	sms.Uid = uid

	if err := VailMobile(sms.Mobile); err != nil {
		return err
	}

	if err := VailUid(sms.Uid); err != nil {
		return err
	}

	sms.model.SetSmsUid()

	return nil
}

func (sms *SMS) DelUid(mobile, uid string) error {
	if !sms.ConfigisOK {
		return fmt.Errorf("(%s)服务配置不存在", sms.serviceName)
	}

	sms.Mobile = mobile
	sms.Uid = uid

	if err := VailMobile(sms.Mobile); err != nil {
		return err
	}
	if err := VailUid(sms.Uid); err != nil {
		return err
	}

	olduid, err := sms.model.GetSmsUid()

	if err != nil {
		return fmt.Errorf(config.Errormsg["err_not_uid"], sms.Mobile, sms.Uid)
	}
	if olduid != uid {
		return fmt.Errorf(config.Errormsg["err_not_uid"], sms.Mobile, sms.Uid)
	}

	sms.model.DelSmsUid()
	return nil
}

func (sms *SMS) Info(mobile string) (map[string]interface{}, error) {
	if !sms.ConfigisOK {
		return nil, fmt.Errorf("(%s)服务配置不存在", sms.serviceName)
	}
	sms.Mobile = mobile

	info := make(map[string]interface{})
	info["mobile"] = sms.Mobile
	info["service"] = sms.serviceName
	info["areacode"], _ = sms.model.GetMobileArea()
	info["lastsendtime"], _ = sms.model.GetSendTime()
	info["sendnums"], _ = sms.model.GetTodaySendNums()
	info["smscode"], info["smscodevalidtime"], _ = sms.model.GetSmsCode()

	info["uid"], _ = sms.model.GetSmsUid()
	return info, nil
}
