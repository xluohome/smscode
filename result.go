package main

type Result interface {
	//返回对象格式化为指定样式
	Format(err error, info interface{})
}

//###########默认#########

type Result_default struct {
	Success bool        `json:"success"`
	Info    interface{} `json:"info"`
}

func (res *Result_default) Format(err error, info interface{}) {

	if err == nil {
		res.Success = true
		res.Info = info
	} else {
		res.Info = err.Error()
	}
	return
}

//###########mobcent#########

type mobcent_head struct {
	ErrCode string `json:"errCode"`
	ErrInfo string `json:"errInfo"`
	Version string `json:"version"`
	Alert   uint8  `json:"alert"`
}

type Result_mobcent struct {
	Rs      int                    `json:"rs"`
	Errcode string                 `json:"errcode"`
	Body    map[string]interface{} `json:"body"`
	Head    mobcent_head           `json:"head"`
}

func (res *Result_mobcent) Format(err error, info interface{}) {
	res.Body = make(map[string]interface{})
	if err == nil {
		res.Rs = 1
		res.Head.ErrCode = "000000"
		res.Head.ErrInfo = "golang调用成功,没有任何错"
		if info != nil {
			res.Body["info"] = info
		}
	} else {
		res.Head.ErrCode = "900000"
		res.Head.ErrInfo = err.Error()
		res.Errcode = err.Error()
	}
	res.Head.Version = "2.7.0.3"
	return
}
