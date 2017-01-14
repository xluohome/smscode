package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const ContentType = "text/json"

type apiserver struct {
}

func (s *apiserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer fmt.Println(r.RemoteAddr, r.URL)

	var err error
	var infos interface{}

	mobile := r.FormValue("mobile")
	code := r.FormValue("code")
	serviceName := r.FormValue("service")
	uid := r.FormValue("uid")

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

	s.echoJson(w, err, infos)
}

func (s *apiserver) echoJson(w http.ResponseWriter, err error, info interface{}) {

	var result Result

	switch sms.Config.Outformat {
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
	fmt.Println("Start Smscode Server...")

	http.ListenAndServe(config.Bind, &apiserver{})
}
