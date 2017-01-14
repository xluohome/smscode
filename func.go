package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

//生成验证码
func makeCode() (code string) {
	code = strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(899999) + 100000)
	return
}

//int64转[]byte
func Int64ToBytes(n int64) (b []byte) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, n)
	return bytesBuffer.Bytes()
}

//[]byte转换int64
func BytesToInt64(b []byte) (tmp int64) {
	bytesBuffer := bytes.NewBuffer(b)
	binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	return
}

//uint64转[]byte
func Uint64ToBytes(n uint64) (b []byte) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, n)
	return bytesBuffer.Bytes()
}

//[]byte转换uint64
func BytesTouInt64(b []byte) (tmp uint64) {
	bytesBuffer := bytes.NewBuffer(b)
	binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	return
}

//[]byte转换uint8
func BytesTouInt8(b []byte) (tmp uint8) {
	bytesBuffer := bytes.NewBuffer(b)
	binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	return
}

func VailMobile(mobile string) error {

	if len(mobile) < 11 {
		return errors.New("[mobile]参数不对")
	}
	reg, err := regexp.Compile("^1[3-8][0-9]{9}$")
	if err != nil {
		panic("regexp error")
	}
	if !reg.MatchString(mobile) {
		return errors.New("手机号码[mobile]格式不正确")
	}
	return nil
}

func VailCode(code string) error {

	if len(code) != 6 {
		return errors.New("[code]参数不对")
	}
	c, err := regexp.Compile("^[0-9]{6}$")
	if err != nil {
		panic("regexp error")
	}
	if !c.MatchString(code) {
		return errors.New("验证码[code]格式不正确")
	}
	return nil
}

func VailUid(uid string) error {

	if len(uid) < 1 || len(uid) > 64 {
		return errors.New("[uid]参数不对")
	}
	return nil
}
