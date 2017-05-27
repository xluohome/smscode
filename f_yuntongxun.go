package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/golang/glog"
)

type Yuntongxun struct {
	sms *SMS
}

func init() {
	SenderMap["yuntongxun"] = func() Sender {
		return &Yuntongxun{}
	}
}

type formdata struct {
	To         string   `json:"to"`
	TemplateId string   `json:"templateId"`
	AppId      string   `json:"appId"`
	Datas      []string `json:"datas"`
}

type result struct {
	StatusCode string `json:"statusCode"`
	StatusMsg  string `json:"statusMsg"`
}

func (y *Yuntongxun) Send(sms *SMS) error {
	y.sms = sms
	b, err := y.body()
	if err != nil {
		return err
	}
	body := bytes.NewReader(b)
	req, err := http.NewRequest("POST", y.url(), body)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("charset", "utf-8")
	req.Header.Add("Authorization", y.authen())

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	res := &result{}
	if err := json.Unmarshal(data, res); err != nil {
		return err
	}

	if res.StatusCode != "000000" {
		log.Errorf("Yuntongxun:%+v", res)
		return fmt.Errorf("%s", res.StatusMsg)
	}
	return nil
}

func (y *Yuntongxun) url() string {
	var vendor = config.Vendors["yuntongxun"]
	return fmt.Sprintf("%s/%s/Accounts/%s/SMS/TemplateSMS?sig=%s", vendor["RestURL"], vendor["SoftVersion"], vendor["AccountSid"], y.sig())
}

func (y *Yuntongxun) sig() string {
	var buf bytes.Buffer
	buf.WriteString(config.Vendors["yuntongxun"]["AccountSid"])
	buf.WriteString(config.Vendors["yuntongxun"]["AccountToken"])
	buf.WriteString(y.sms.NowTime.Format("20060102150405"))

	var md5hex = md5.New()
	md5hex.Write(buf.Bytes())

	var sig = strings.ToUpper(hex.EncodeToString(md5hex.Sum(nil)))
	return sig
}

func (y *Yuntongxun) body() ([]byte, error) {
	var datas = []string{y.sms.Code, fmt.Sprintf("%d", y.sms.Config.Validtime/60)} //单位是分钟
	var fd = formdata{y.sms.Mobile, y.sms.Config.Tpl, config.Vendors["yuntongxun"]["AppId"], datas}

	return json.Marshal(fd)
}

// 生成授权：主帐户Id + 英文冒号 + 时间戳。
func (y *Yuntongxun) authen() string {
	var buf bytes.Buffer
	buf.WriteString(config.Vendors["yuntongxun"]["AccountSid"])
	buf.WriteByte(':')
	buf.WriteString(y.sms.NowTime.Format("20060102150405"))
	return base64.URLEncoding.EncodeToString(buf.Bytes())
}
