// Copyright 2021-2022 The jdh99 Authors. All rights reserved.
// 网络校时服务
// Authors: jdh99 <jdh821@163.com>

package main

import (
	"github.com/jdhxyy/arrow"
	"github.com/jdhxyy/dcom"
	"github.com/jdhxyy/lagan"
	"github.com/jdhxyy/utz"
	"ntp/config"
	"time"
)

const tag = "ntp"

// rid号
const (
	// 读取时间.返回的是字符串
	ridGetTime1 = 1
	// 读取时间.返回的是结构体
	ridGetTime2 = 2
)

// AckRidGetTime2 ACK格式
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
	lagan.SetFilterLevel(lagan.LevelDebug)

	err = arrow.Load(config.LocalIA, config.LocalIP, config.LocalPort, config.CoreIA, config.CoreIP, config.CorePort)
	if err != nil {
		lagan.Error(tag, "arrow load failed!")
		return
	}

	arrow.Register(utz.HeaderCcp, ridGetTime1, ntpService1)
	arrow.Register(utz.HeaderCcp, ridGetTime2, ntpService2)

	select {}
}

func ntpService1(req []uint8, params ...interface{}) []uint8 {
	ia := params[0].(uint32)
	ip := params[1].(uint32)
	port := params[2].(uint16)

	var timeZone int
	if len(req) == 0 {
		timeZone = 8
	} else if len(req) == 1 {
		timeZone = int(int8(req[0]))
	} else {
		lagan.Warn(tag, "addr:0x%08x:%d ia:0x%x ntp failed.len is wrong:%d", ip, port, ia, len(req))
		return nil
	}

	t := getTime(timeZone)

	lagan.Info(tag, "addr:0x%08x:%d ia:0x%x ntp time:%v", ip, port, ia, t)
	return []uint8(t.Format("2006-01-02 15:04:05 -0700 MST"))
}

func getTime(timeZone int) time.Time {
	t := time.Now().UTC()
	secondsEastOfUTC := int((time.Duration(timeZone) * time.Hour).Seconds())
	loc := time.FixedZone("CST", secondsEastOfUTC)
	t = t.In(loc)
	return t
}

func ntpService2(req []uint8, params ...interface{}) []uint8 {
	ia := params[0].(uint32)
	ip := params[1].(uint32)
	port := params[2].(uint16)

	var timeZone int
	if len(req) == 0 {
		timeZone = 8
	} else if len(req) == 1 {
		timeZone = int(int8(req[0]))
	} else {
		lagan.Warn(tag, "addr:0x%08x:%d ia:0x%08x ntp failed.len is wrong:%d", ip, port, ia, len(req))
		return nil
	}

	t := getTime(timeZone)
	lagan.Info(tag, "addr:0x%08x:%d ia:0x%x ntp time:%v", ip, port, ia, t)

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
		lagan.Error(tag, "addr:0x%08x:%d ia:0x%x ntp failed.struct to bytes error:%v", ip, port, ia, err)
		return nil
	}
	return data
}
