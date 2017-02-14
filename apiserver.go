package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/golang/glog"
)

type apiserver struct {
}

func (apiserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer log.Info(r.RemoteAddr, r.URL)

	var err error
	var infos interface{}

	mobile := r.FormValue("mobile")
	code := r.FormValue("code")
	serviceName := r.FormValue("service")
	uid := r.FormValue("uid")

	sms := NewSms()
	sms.SetServiceConfig(serviceName)

	switch r.URL.String() {
	case "/send":
		err = sms.Send(mobile)
	case "/checkcode":
		err = sms.CheckCode(mobile, code)
	case "/setuid":
		err = sms.SetUid(mobile, uid)
	case "/deluid":
		err = sms.DelUid(mobile, uid)
	case "/info":
		infos, err = sms.Info(mobile)
	default:
		err = fmt.Errorf("%s", "您访问的api不存在")
	}

	var result Result

	switch sms.Config.Outformat {
	case "mobcent":
		result = &Result_mobcent{}
	default:
		result = &Result_default{}
	}

	//对象输出格式化
	result.Format(err, infos)

	//json对象
	str, _ := json.Marshal(result)
	w.Header().Set("Content-Type", ContentType)
	w.Write(str)
}

func Apiserver() {
	http.ListenAndServe(config.Bind, &apiserver{})
}
