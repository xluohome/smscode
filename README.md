SmsCode
=======

### 这是什么

  专门为需要【手机短信验证码】使用场景而设计的微服务(Micro service)，如：用户注册、找回密码、用户身份验证、验证码登录等等。

### 为什么要造这个轮子

首先我没有找到已实现这些基本功能合一的轮子，至少开源的没有。而实际上我们很多项目开发时经常需要用到手机验证码功能。然而每次重复造轮子又觉得太繁琐且不容易集中控制。于是有了开发一个可复用轮子的想法。那么合计不到1500行Go代码实现这个可复用的手机验证码微服务何乐不为呢？

### 安装使用

得益于Go语言的跨平台支持，SmsCode可安装在所有主流OS上（Linux，Mac OS X，FreeBSD，Windows，ARM等）

推荐Linux x64上安装SmsCode，编译安装请确保已经在OS上安装了Go的编译环境(GO 1.5+)。

    go get -u github.com/xluohome/smscode

	cd  $GOPATH/src/github.com/xluohome/smscode

	./build  && ./smscode

### Docker部署

请参考项目中的Dockerfile 制作Docker image。

### 时区设置

请修改 conf.yaml 中的 timezone 参数。

timezone 具体值请查考 conf/zoneinfo.zip

### 功能特性

1. 支持阿里大鱼、云通讯等多个“手机短信验证码”(以下简称：验证码)通道;
2. 自定义多个验证码服务接口，如：新用户注册，重设密码，身份验证等等;
3. 内置手机号归属地限制，限制指定归属地手机号可接收验证码;
4. 每个验证码服务可设置“验证码”每日最大发送数量限额及有效时间;
5. 内置callback服务，可设置验证码发送成功（失败）、验证码验证成功时的回调http Url;
6. 可设置验证码发送模式:
 - 0x01：只有手机号对应的uid存在时才能发送。
 - 0x02：只有uid不存在时才能发送。
 - 0x03：不管uid是否存在都发送。
7. 通过setuid接口可将现有系统中的用户UID数据导入SmsCode;
8. 内置本地持久化存储：Goleveldb;
9. 支持Docker部署，SmsCode静态编译(Go 1.7.5)Docker image不到12mb(含归属地数据库);

### 配置文件 conf/conf.yaml

```
bind: 0.0.0.0:8080  #短信验证码微服务器地址
timezone: PRC   #时区设置
timeout: 5 #短信供应商网关响应超时时间（秒）
vendors:
  alidayu: #阿里大鱼配置 http://www.alidayu.com
    appkey: 20315570
    appSecret: 87hfgfg75775765787878
    issendbox: false
  yuntongxun: #云通讯配置  http://www.yuntongxun.com/
    AccountSid: 8a48b55e434514c9c31921a039b
    AccountToken: 61434dc2b245435eadf82d381fa3f
    AppId: aaf98f8fsdafd2678c9d07875040f
    SoftVersion: 2013-12-26
    RestURL: https://app.cloopen.com:8883
  hywx: #互亿无线  http://www.ihuyi.cn/
    account: 6666666666666666
    password: 88888888888888888888888
    RestURL: http://106.ihuyi.cn/webservice/sms.php?method=Submit

errormsg:
  "err_model_not_ok1": "当前用户(%s)不存在，不能发送手机验证码"
  "err_model_not_ok2": "当前用户(%s)存在，不能发送手机验证码"
  "err_code_not_ok": "手机验证码:%s 不正确，请重新输入"
  "err_vailtime_not_ok": "手机验证码已超时:%s，请重新获取"
  "err_per_day_max_send_nums": "一个手机号每天仅限发送%d条验证码"
  "err_per_minute_send_num": "一个手机号每分钟仅限发送1条验证码"
  "err_allow_areacode": "手机号归属地不允许，仅限于:%s"
  "err_not_uid": "手机号%s对应的%s不存在或者不匹配"

servicelist:
  "register":
    vendor: alidayu  #短信通道供应商
    group: db1   #相同组内的uid数据共享
    smstpl: SMS_34850248  #阿里大鱼短信模板id
    signname: 罗永 #阿里大鱼短信签名
    callback: "http://127.0.01/test9.php"
    allowcitys: #仅限如下的手机号归属区接收验证码
      - 0575
      - 0571
      - 0574
    maxsendnums: 4   #一个手机号每天发送限额,这个受短信运营商的限制。
    validtime: 600  #单位：秒 。 收到的手机验证码x秒内有效，超过后验证无效；
    mode: 2   #模式  1：只有手机号对应的uid存在时才能发送，2：只有uid不存在时才能发送，3：不管uid是否存在都发送
    outformat: mobcent  #RestAPi接口输出样式（mobcent,default）

  "restpwd":
    vendor: alidayu
    group: db1
    smstpl: SMS_39190087
    signname: 罗永
    callback:
    allowcitys:
      - 0578
    maxsendnums: 2
    validtime: 360
    mode: 3

  "getpwd":
    vendor: yuntongxun
    group: db1
    smstpl: 149350
    signname: 罗永亿
    callback:
    allowcitys:
      - 0578
      - 0575
    maxsendnums: 2
    validtime: 360
    mode: 3
```
### SmsCode RestAPi 接口说明


| 接口名称        | 接口地址           | 参数  | 说明 |
  :------------- |:-------------| :-----|:----|
|发送验证码 |/send	| service, mobile   |服务名称,手机号|
|验证验证码	|/checkcode	|service, mobile,code   |服务名称,手机号,验证码
|设置用户UID |/setuid | service, mobile,uid |服务名称,手机号 ,用户id
|删除用户UID	|/deluid |service, mobile,uid |服务名称,手机号,用户id
|信息查询|/info |service, mobile  |服务名称,手机号

### SmsCode RestAPi 接口返回说明
具体依据返回内容的结构体。
以下是 mobcent结构体json内容

#### 执行成功：
```
{
  "rs": 1,
  "errcode": "",
  "body": {},
  "head": {
    "errCode": "000000",
    "errInfo": "smscode调用成功,没有任何错",
    "version": "2.7.0.3",
    "alert": 0
  }
}
```
#### 执行失败：
```
{
  "rs": 0,
  "errcode": "[mobile]参数不对",
  "body": {},
  "head": {
    "errCode": "900000",
    "errInfo": "[mobile]参数不对",
    "version": "2.7.0.3",
    "alert": 0
  }
}
```
#### 信息查询返回：
```
{
  "rs": 1,
  "errcode": "",
  "body": {
    "info": {
      "areacode": "0575",
      "extinfo": {
        "PhoneNum": "15336860197",
        "Province": "浙江",
        "City": "绍兴",
        "ZipCode": "312000",
        "AreaZone": "0575",
        "CardType": "电信"
      },
      "lastsendtime": 0,
      "mobile": "15336860197",
      "sendnums": 0,
      "service": "register",
      "smscode": "",
      "smscodeinvalidtime": 0,
      "uid": ""
    }
  },
  "head": {
    "errCode": "000000",
    "errInfo": "smscode调用成功,没有任何错",
    "version": "2.7.0.3",
    "alert": 0
  }
}
```

### Callback 服务

每个短信验证码服务允许设置一个callback http Url ；
如下事件发生时将回调callback http Url

1. 手机验证码发送成功或者失败时。
2. 手机验证码通过验证时。

当 callback http Url无法访问时，系统会延时2,4,6,8,16,32,64,128秒 依次进行重试（合计10次）。

成功时附带如下Url POST请求参数(multipart/form-data)：
```

mobile string  #手机号
code string	#验证码
service string #服务名称
uxtime int #时间戳
flag  string #回调标记 (Success,Failed,Checkok)
````

### 联系作者
欢迎来信交流
phposs@qq.com ，
QQ群
309020981
