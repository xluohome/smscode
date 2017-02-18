package main

import (
	"net/http"
	"testing"
)

func TestNewCallback(t *testing.T) {

	wating := make(chan bool)

	var fun = func(w http.ResponseWriter, r *http.Request) {

		t.Log("mobile:", r.FormValue("mobile"), "code:", r.FormValue("code"), "uxtime:", r.FormValue("uxtime"))
		w.Write([]byte("ok"))
		wating <- true
	}

	go func() {
		http.HandleFunc("/test", fun)
		http.ListenAndServe("127.0.0.1:8080", nil)

	}()

	sms := NewSms()
	sms.SetServiceConfig("register")
	sms.Mobile = "13575566313"
	sms.Code = "999999"
	sms.Config.Callback = "http://127.0.0.1:8080/test"
	AddCallbackTask(sms, "test")

	<-wating

}
