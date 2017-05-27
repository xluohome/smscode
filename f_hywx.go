/**
互亿无线 短信通道
http://www.ihuyi.cn/
**/
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/golang/glog"
)

type Hywx struct {
	sms     *SMS
	XMLName xml.Name `xml:"SubmitResult"`
	Code    int      `xml:"code"`
	Msg     string   `xml:"msg"`
}

func init() {
	SenderMap["hywx"] = func() Sender {
		return &Hywx{}
	}
}

func (h *Hywx) Send(sms *SMS) error {
	h.sms = sms
	var data = make(url.Values)
	data.Set("account", config.Vendors["hywx"]["account"])
	data.Set("password", config.Vendors["hywx"]["password"])
	data.Set("mobile", h.sms.Mobile)
	data.Set("content", strings.Replace(h.sms.Config.Tpl, "{code}", h.sms.Code, -1))
	res, err := http.PostForm(config.Vendors["hywx"]["RestURL"], data)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := xml.Unmarshal(body, h); err != nil {
		return err
	}

	if h.Code != 2 {
		log.Errorf("%v", h)
		return fmt.Errorf("%s", h.Msg)
	}

	return nil
}
