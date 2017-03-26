package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/golang/glog"
)

type apiserver struct {
	req    []chan *request
	result []chan Result
}

type request struct {
	Act, Mobile, Code, Service, Uid string `json:",omitempty"`
}

func (r *request) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s", r.Act, r.Mobile, r.Code, r.Service, r.Uid)
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
		api.result[i] <- result
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

	req := request{r.URL.String()[1:], r.FormValue("mobile"), r.FormValue("code"), r.FormValue("service"), r.FormValue("uid")}
	hash := hashFunc([]byte(req.String())) & uint64(*smsworks-1)
	api.req[hash] <- &req

	var result Result
	select {
	case result = <-api.result[hash]:
	case <-time.After(config.TimeOut * time.Second):
		panic("Response timeout error!")
	}
	w.Header().Set("Content-Type", ContentType)
	w.Header().Set("TimeZone", time.Local.String())
	str, _ := json.Marshal(result)
	w.Write(str)
}

func Apiserver() {
	var api = new(apiserver)
	api.req = make([]chan *request, *smsworks)
	api.result = make([]chan Result, *smsworks)
	for i := 0; i < *smsworks; i++ {
		go func(i int) {
			api.req[i] = make(chan *request)
			api.result[i] = make(chan Result)
			api.sms(i)
		}(i)
	}

	log.Fatal(http.ListenAndServe(config.Bind, api))
}
