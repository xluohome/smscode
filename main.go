package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

var (
	config          Config
	configFile      = flag.String("c", "conf/conf.yaml", "配置文件,default:conf/conf.yaml")
	dbPath          = flag.String("db", "level.db", "数据库保存路径")
	smsworks        = flag.Int("sw", runtime.NumCPU(), "短信验证码服务器数量,default:cpu number")
	callbackworkers = flag.Int("cw", runtime.NumCPU(), "Callback并发数,default:cpu number")
	trycallnums     = flag.Uint("cn", 10, "Callback失败重试次数")
)

func init() {

	flag.Usage = func() {
		fmt.Printf("Usage %s -c=conf/conf.yaml -db level.db -sw 10 -cw 10 -cn 10", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := config.ParseConfigFile(*configFile); err != nil {
		fmt.Printf("%s\n%s", "Unable to load config info!", err)
		os.Exit(1)
	}

	time.Local = func() *time.Location {
		loc, err := time.LoadLocation(config.Timezone)
		if err != nil {
			loc = time.UTC
		}
		return loc
	}()

	RunCallbackTask()
}

func main() {
	smscodetxt := `
                                             **
                                              *
                                              *
  ***** *******   *****    ***    ****    *****   ****
 *    *  *  *  * *    *   *   *  *    *  *    *  *    *
 *       *  *  * *       *       *    *  *    *  *    *
  ****   *  *  *  ****   *       *    *  *    *  ******
      *  *  *  *      *  *       *    *  *    *  *
 *    *  *  *  * *    *   *   *  *    *  *   **  *    *
 *****  *** ** *******     ***    ****    *** **  ****
https://github.com/xluohome/smscode ©xluo(phposs@qq.com)

`
	fmt.Print(smscodetxt)
	fmt.Printf("Smscode Server v%s  ...    TimeZone:%s \n\n", VERSION, config.Timezone)
	Apiserver()
}
