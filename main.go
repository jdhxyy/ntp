// Copyright 2021-2021 The jdh99 Authors. All rights reserved.
// 网络校时服务
// Authors: jdh99 <jdh821@163.com>

package main

import (
	"github.com/jdhxyy/dcom"
	"github.com/jdhxyy/lagan"
	"github.com/jdhxyy/tziot"
	"ntp/config"
	"time"
)

const tag = "ntp"

// 应用错误码
const (
	// 内部错误
	errorCodeInternalError = 0x40
	// 接收格式错误
	errorCodeRxFormat = 0x41
)

// rid号
const (
	// 读取时间.返回的是字符串
	ridGetTime1 = 1
	// 读取时间.返回的是结构体
	ridGetTime2 = 2
)

// ACK格式
type AckRidGetTime2 struct {
	// 时区
	TimeZone uint8
	Year     uint16
	Month    uint8
	Day      uint8
	Hour     uint8
	Minute   uint8
	Second   uint8
	// 星期
	Weekday uint8
}

func main() {
	err := lagan.Load(0)
	if err != nil {
		panic(err)
	}
	lagan.EnableColor(true)
	lagan.SetFilterLevel(lagan.LevelInfo)

	pipe := tziot.BindPipeNet(config.LocalIA, config.LocalPwd, config.LocalIP, config.LocalPort)
	if pipe == 0 {
		lagan.Error(tag, "bind pipe failed!")
		return
	}
	tziot.Register(ridGetTime1, ntpService1)
	tziot.Register(ridGetTime2, ntpService2)

	select {}
}

// ntpService1 校时服务
// 返回值是应答和错误码.错误码为0表示回调成功,否则是错误码
func ntpService1(pipe uint64, srcIA uint64, req []uint8) ([]uint8, int) {
	addr := dcom.PipeToAddr(pipe)

	var timeZone int
	if len(req) == 0 {
		timeZone = 8
	} else if len(req) == 1 {
		timeZone = int(int8(req[0]))
	} else {
		lagan.Warn(tag, "addr:%v ia:0x%x ntp failed.len is wrong:%d", addr, srcIA, len(req))
		return nil, errorCodeRxFormat
	}

	t := getTime(timeZone)
	lagan.Info(tag, "addr:%v ia:0x%x ntp time:%v", addr, srcIA, t)
	return []uint8(t.Format("2006-01-02 15:04:05 -0700 MST")), 0
}

func getTime(timeZone int) time.Time {
	t := time.Now().UTC()
	secondsEastOfUTC := int((time.Duration(timeZone) * time.Hour).Seconds())
	loc := time.FixedZone("CST", secondsEastOfUTC)
	t = t.In(loc)
	return t
}

// ntpService2 校时服务
// 返回值是应答和错误码.错误码为0表示回调成功,否则是错误码
func ntpService2(pipe uint64, srcIA uint64, req []uint8) ([]uint8, int) {
	addr := dcom.PipeToAddr(pipe)

	var timeZone int
	if len(req) == 0 {
		timeZone = 8
	} else if len(req) == 1 {
		timeZone = int(int8(req[0]))
	} else {
		lagan.Warn(tag, "addr:%v ia:0x%x ntp failed.len is wrong:%d", addr, srcIA, len(req))
		return nil, errorCodeRxFormat
	}

	t := getTime(timeZone)
	lagan.Info(tag, "addr:%v ia:0x%x ntp time:%v", addr, srcIA, t)

	var ack AckRidGetTime2
	ack.TimeZone = uint8(timeZone)
	ack.Year = uint16(t.Year())
	ack.Month = uint8(t.Month())
	ack.Day = uint8(t.Day())
	ack.Hour = uint8(t.Hour())
	ack.Minute = uint8(t.Minute())
	ack.Second = uint8(t.Second())
	ack.Weekday = uint8(t.Weekday())

	data, err := dcom.StructToBytes(ack)
	if err != nil {
		lagan.Error(tag, "addr:%v ia:0x%x ntp failed.struct to bytes error:%v", addr, srcIA, err)
		return nil, errorCodeInternalError
	}
	return data, 0
}
