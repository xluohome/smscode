package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/golang/glog"
)

type JuheApi struct {
	Url    string
	Phone  string
	Key    string
	Result *ApiResult
}

type ApiResult struct {
	Resultcode string            `json:"resultcode"`
	Result     map[string]string `json:"result"`
	Error_code int               `json:"error_code"`
}

func NewJuheApi() *JuheApi {

	return &JuheApi{Url: "http://apis.juhe.cn/mobile/get", Key: config.Juheapikey}

}

func (j *JuheApi) Query(phone string) error {
	var url = fmt.Sprintf("%s?phone=%s&key=%s", j.Url, phone, j.Key)
	log.V(1).Info(url)
	res, err := http.Get(url)
	if err != nil {
		log.Error("JuheApi访问出错！！！")
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.V(1).Infof("%s", body)
	if err := json.Unmarshal(body, &j.Result); err != nil {
		return err
	}
	return nil
}
