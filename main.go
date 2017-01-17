package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	config          Config
	configFile      = flag.String("c", "etc/conf.yaml", "配置文件，默认etc/conf.yaml")
	dbPath          = flag.String("db", "level.db", "database path")
	callbackworkers = flag.Int("cw", runtime.NumCPU(), "callback并发数，默认是cpu数")
	trycallnums     = flag.Int("cn", 10, "callback失败重试次数")
	sms             *SMS
	SMSModel        *Model
)

func init() {

	flag.Usage = func() {
		fmt.Printf("Usage %s -c=etc/conf.yaml -db level.db -cw 10 -cn 10", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := config.ParseConfigFile(*configFile); err != nil {
		fmt.Printf("%s,%s", "配置文件加载错误", err)
		os.Exit(1)
	}

	sms = NewSms()

	SMSModel = NewModel(sms)

	RunCallbackTask()
}

func main() {

	Apiserver()
}
