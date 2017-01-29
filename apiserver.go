package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiserver struct {
	sms *SMS
}

func (s *apiserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer fmt.Println(r.RemoteAddr, r.URL)

	var err error
	var infos interface{}

	mobile := r.FormValue("mobile")
	code := r.FormValue("code")
	serviceName := r.FormValue("service")
	uid := r.FormValue("uid")

	s.sms = NewSms()
	s.sms.SetServiceConfig(serviceName)

	switch r.URL.String() {
	case "/send":
		err = s.sms.Send(mobile)
	case "/checkcode":
		err = s.sms.CheckCode(mobile, code)
	case "/setuid":
		err = s.sms.SetUid(mobile, uid)
	case "/deluid":
		err = s.sms.DelUid(mobile, uid)
	case "/info":
		infos, err = s.sms.Info(mobile)
	default:
		err = fmt.Errorf("%s", "您访问的api不存在")
	}

	s.echoJson(w, err, infos)
}

func (s *apiserver) echoJson(w http.ResponseWriter, err error, info interface{}) {

	var result Result

	switch s.sms.Config.Outformat {
	case "mobcent":
		result = &Result_mobcent{}
	default:
		result = &Result_default{}
	}

	//对象输出格式化
	result.Format(err, info)

	//json对象
	str, _ := json.Marshal(result)
	w.Header().Set("Content-Type", ContentType)
	w.Write(str)
}

func Apiserver() {
	http.ListenAndServe(config.Bind, &apiserver{})
}
