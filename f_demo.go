//测试用 短信通道供应商
package main

import (
	"fmt"
	"time"
)

type Demov struct {
	sms *SMS
}

func (y *Demov) Send(sms *SMS) error {

	time.Sleep(5 * time.Second)
	fmt.Printf("发送成功")
	return nil
}
