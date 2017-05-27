package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/golang/glog"
)

type request struct {
	Act, Mobile, Code, Service, Uid string `json:",omitempty"`
	result                          chan Result
}

func (r *request) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s", r.Act, r.Mobile, r.Code, r.Service, r.Uid)
}

type apiserver struct {
	req []chan *request
}

func (api *apiserver) sms(i int) {

	var (
		infos  interface{}
		err    error
		result Result
	)
	sms := NewSms()
	for req := range api.req[i] {

		sms.NowTime = time.Now()
		sms.SetServiceConfig(req.Service)

		infos = nil

		switch req.Act {
		case "send":
			err = sms.Send(req.Mobile)
		case "checkcode":
			err = sms.CheckCode(req.Mobile, req.Code)
		case "setuid":
			err = sms.SetUid(req.Mobile, req.Uid)
		case "deluid":
			err = sms.DelUid(req.Mobile, req.Uid)
		case "info":
			infos, err = sms.Info(req.Mobile)
		default:
			err = fmt.Errorf("%s\t%s", "API not found:", req.Act)
		}

		switch sms.Config.Outformat {
		case "mobcent":
			result = &Result_mobcent{}
		default:
			result = &Result_default{}
		}

		result.Format(err, infos)
		req.result <- result
	}
	return
}

func (api *apiserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info(r.RemoteAddr, r.URL)
	defer func() {
		if e := recover(); e != nil {
			err := fmt.Errorf("Server error :%s", e)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
	}()

	resultchan := make(chan Result)
	req := &request{r.URL.String()[1:], r.FormValue("mobile"), r.FormValue("code"), r.FormValue("service"), r.FormValue("uid"), resultchan}
	hash := hashFunc([]byte(req.String())) & uint64(*smsworks-1)

	for {
		select {
		case api.req[hash] <- req:
		case result := <-resultchan:
			w.Header().Set("Content-Type", ContentType)
			w.Header().Set("TimeZone", time.Local.String())
			str, _ := json.Marshal(result)
			w.Write(str)
			return

			//未在预定时间内完成请求或者接收答复则报超时错误
			//如果设置的smsworks过少时，当大量的并发请求因处理量受限无法及时处理。会使得请求或接受时间延长导致超时；
			//另外网络不稳定或短信供应商通道出现了问题也会引起超时。
		case <-time.After(config.TimeOut * time.Second):
			panic("Server timeout error!")
		}
	}
}

func Apiserver() {

	var api = new(apiserver)
	api.req = make([]chan *request, *smsworks)
	for i := 0; i < *smsworks; i++ {
		go func(i int) {
			api.req[i] = make(chan *request)
			api.sms(i)
		}(i)
	}

	log.Fatal(http.ListenAndServe(config.Bind, api))
}
