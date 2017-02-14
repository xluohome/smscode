package main

import (
	"bytes"
	//"fmt"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/zheng-ji/gophone"
)

type Model struct {
	sms *SMS
	db  *leveldb.DB
	mu  sync.Mutex
}

var (
	db   *leveldb.DB
	once sync.Once
)

func init_db() {
	var err error
	var opt = &opt.Options{BlockCacheCapacity: CacheSize,
		WriteBuffer: WriteBuffer * 1024 * 1024}
	db, err = leveldb.OpenFile(*dbPath, opt)
	if err != nil {
		log.Fatalln("db.Get(), err:", err)
	}
}

func NewModel(sms *SMS) (SMSModel *Model) {

	once.Do(func() {
		init_db()
	})
	SMSModel = &Model{}
	SMSModel.sms = sms
	SMSModel.db = db
	return
}

func (m *Model) GetMobileInfo() (*gophone.PhoneRecord, error) {

	//感谢 zheng-ji github.com/zheng-ji/gophone
	pr, err := gophone.Find(m.sms.Mobile)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (m *Model) GetMobileArea() (string, error) {

	//感谢 zheng-ji github.com/zheng-ji/gophone
	mobileinfo, err := m.GetMobileInfo()
	if err != nil {
		return "", err
	}
	return mobileinfo.AreaZone, nil
}

/**
获取发送时间
group:signname:mobile
**/
func (m *Model) GetSendTime() (int64, error) {

	var buf bytes.Buffer
	buf.WriteString(m.sms.Config.Group)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Config.Signname)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Mobile)
	buf.WriteString(":sendTime:")

	key := buf.Bytes()
	sendTime, err := m.db.Get(key, nil)
	if err == errors.ErrNotFound {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return BytesToInt64(sendTime), nil
}

/**
设置发送时间
group:signname:mobile
**/
func (m *Model) SetSendTime() {

	var buf bytes.Buffer
	buf.WriteString(m.sms.Config.Group)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Config.Signname)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Mobile)
	buf.WriteString(":sendTime:")

	key := buf.Bytes()
	time := m.sms.NowTime.Unix()
	val := []byte(Int64ToBytes(time))

	m.db.Put(key, val, nil)
}

/**
获取当天发送数量
group:service:signname:mobile:date
**/
func (m *Model) GetTodaySendNums() (uint64, error) {

	var buf bytes.Buffer
	buf.WriteString(m.sms.Config.Group)
	buf.WriteByte(':')
	buf.WriteString(m.sms.serviceName)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Config.Signname)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Mobile)
	buf.WriteByte(':')
	buf.WriteString(m.sms.NowTime.Format("2006-01-02"))
	buf.WriteString(":sendNums:")

	key := buf.Bytes()
	nums, err := m.db.Get(key, nil)
	if err == errors.ErrNotFound {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return BytesTouInt64(nums), nil
}

/**
设置当天发送数量
group:service:signname:mobile:date
**/
func (m *Model) SetTodaySendNums(num uint64) {

	var buf bytes.Buffer
	buf.WriteString(m.sms.Config.Group)
	buf.WriteByte(':')
	buf.WriteString(m.sms.serviceName)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Config.Signname)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Mobile)
	buf.WriteByte(':')
	buf.WriteString(m.sms.NowTime.Format("2006-01-02"))
	buf.WriteString(":sendNums:")

	key := buf.Bytes()
	val := Uint64ToBytes(num)

	m.db.Put(key, val, nil)
}

/**
获取code
group:servicename:signname:mobile
**/
func (m *Model) GetSmsCode() (code string, uxtime int64, err error) {

	var buf bytes.Buffer
	buf.WriteString(m.sms.Config.Group)
	buf.WriteByte(':')
	buf.WriteString(m.sms.serviceName)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Config.Signname)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Mobile)
	buf.WriteString(":SmsCode:")

	key := buf.Bytes()

	var val []byte
	val, err = m.db.Get(key, nil)
	if err == errors.ErrNotFound {
		return "", 0, nil
	} else if err != nil {
		return "", 0, err
	}

	code = string(val[:6])
	if len(val) > 6 {
		uxtime = BytesToInt64(val[7:])
	}

	return
}

/**
设置code
group:servicename:signname:mobile
**/
func (m *Model) SetSmsCode() {

	var buf bytes.Buffer
	buf.WriteString(m.sms.Config.Group)
	buf.WriteByte(':')
	buf.WriteString(m.sms.serviceName)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Config.Signname)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Mobile)
	buf.WriteString(":SmsCode:")

	key := buf.Bytes()

	var val bytes.Buffer
	val.WriteString(m.sms.Code)
	val.WriteByte(':')

	nextvaliedunixtime := m.sms.NowTime.Add(time.Duration(m.sms.Config.Validtime) * time.Second).Unix()

	val.Write([]byte(Int64ToBytes(nextvaliedunixtime)))

	m.db.Put(key, val.Bytes(), nil)
}

/**
获取Uid
group:mobile
**/
func (m *Model) GetSmsUid() (string, error) {

	var buf bytes.Buffer
	buf.WriteString(m.sms.Config.Group)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Mobile)
	buf.WriteString(":SmsUid:")

	key := buf.Bytes()
	uid, err := m.db.Get(key, nil)
	if err == errors.ErrNotFound {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return string(uid), nil
}

/**
设置Uid
group:mobile
**/
func (m *Model) SetSmsUid() {

	var buf bytes.Buffer
	buf.WriteString(m.sms.Config.Group)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Mobile)
	buf.WriteString(":SmsUid:")

	key := buf.Bytes()
	val := []byte(m.sms.Uid)

	m.db.Put(key, val, nil)
}

/**
删除Uid
group:mobile
**/
func (m *Model) DelSmsUid() {

	var buf bytes.Buffer
	buf.WriteString(m.sms.Config.Group)
	buf.WriteByte(':')
	buf.WriteString(m.sms.Mobile)
	buf.WriteString(":SmsUid:")

	key := buf.Bytes()

	m.db.Delete(key, nil)
}
