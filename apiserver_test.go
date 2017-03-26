package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestApiserver(t *testing.T) {

	if e := recover(); e != nil {
		t.Error(e)
	}

	config.Bind = "0.0.0.0:8080"
	go func() {
		Apiserver()
	}()

	t.Run("Send", testsend)
	t.Run("Checkcode", testcheckcode)
	t.Run("Setuid", testsetuid)
	t.Run("Deluid", testdeluid)
	t.Run("Info", testinfo)
}

func testsend(t *testing.T) {
	apiurl := "http://" + config.Bind
	var data = url.Values{}
	data.Set("mobile", "13575566313")
	data.Set("service", "register")
	res, err := http.DefaultClient.PostForm(apiurl+"/send", data)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()
	cont, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(cont))
}

func testcheckcode(t *testing.T) {
	apiurl := "http://" + config.Bind
	var data = url.Values{}
	data.Set("mobile", "13575566313")
	data.Set("service", "register")
	data.Set("code", "333333")
	res, err := http.DefaultClient.PostForm(apiurl+"/checkcode", data)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()
	cont, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(cont))
}

func testsetuid(t *testing.T) {
	apiurl := "http://" + config.Bind
	var data = url.Values{}
	data.Set("mobile", "13575566313")
	data.Set("service", "register")
	data.Set("uid", "1")
	res, err := http.DefaultClient.PostForm(apiurl+"/setuid", data)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()
	cont, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(cont))
}

func testdeluid(t *testing.T) {
	apiurl := "http://" + config.Bind
	var data = url.Values{}
	data.Set("mobile", "13575566313")
	data.Set("service", "register")
	data.Set("uid", "1")
	res, err := http.DefaultClient.PostForm(apiurl+"/deluid", data)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()
	cont, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(cont))
}

func testinfo(t *testing.T) {
	apiurl := "http://" + config.Bind
	var data = url.Values{}
	data.Set("mobile", "13575566313")
	data.Set("service", "register")
	res, err := http.DefaultClient.PostForm(apiurl+"/info", data)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()
	cont, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(cont))
}
