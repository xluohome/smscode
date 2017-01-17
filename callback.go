package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	log "github.com/golang/glog"
)

type Callback struct {
	Url      string
	Data     url.Values
	callnums uint
}

var (
	callback = make(chan Callback)
)

func RunCallbackTask() {
	c := &Callback{}
	for i := 0; i < *callbackworkers; i++ {
		go c.Do(callback)
	}
}

func AddCallbackTask(sms SMS, flag string) {

	if len(sms.Config.Callback) < 1 { //没有启用
		return
	}
	data := make(url.Values)
	data.Set("mobile", sms.Mobile)
	data.Set("code", sms.Code)
	data.Set("service", sms.serviceName)
	data.Set("uxtime", fmt.Sprintf("%d", sms.NowTime.Unix()))
	data.Set("flag", flag)
	callback <- Callback{sms.Config.Callback, data, 0}
}

func (c Callback) Do(cbs <-chan Callback) {

	for cb := range cbs {

		go func() {
			//延时 2,4,6,8,16,32,64,128 ... 秒
			<-time.After(time.Duration(1<<cb.callnums) * time.Second)

			res, err := http.PostForm(cb.Url, cb.Data)
			if err != nil {
				log.Errorf("http.PostForm发生了错误:%v,%s", cb, err.Error())
				return
			}
			defer res.Body.Close()

			//返回200 即callback成功
			if res.StatusCode == http.StatusOK {
				return
			}
			if cb.callnums > trycallnums {
				log.Errorf("重试%d次callback均失败:%v", trycallnums, cb)
				return
			}

			cb.callnums++

			log.V(1).Infof("cb:%v", cb)

			callback <- cb

		}()

	}
}
