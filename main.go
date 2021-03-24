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
	// 接收格式错误
	errorCodeRxFormat = 0x40
)

// rid号
const ridGetTime = 1

func main() {
	err := lagan.Load(0)
	if err != nil {
		panic(err)
	}
	lagan.EnableColor(true)
	lagan.SetFilterLevel(lagan.LevelDebug)

	_, err = tziot.BindPipeNet(config.LocalIA, config.LocalPwd, config.LocalIP, config.LocalPort)
	if err != nil {
		panic(err)
		return
	}
	tziot.Register(ridGetTime, ntpService)

	select {}
}

// ntpService 校时服务
// 返回值是应答和错误码.错误码为0表示回调成功,否则是错误码
func ntpService(pipe uint64, srcIA uint64, req []uint8) ([]uint8, int) {
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

	t := time.Now().UTC()
	secondsEastOfUTC := int((time.Duration(timeZone) * time.Hour).Seconds())
	loc := time.FixedZone("CST", secondsEastOfUTC)
	t = t.In(loc)

	lagan.Info(tag, "addr:%v ia:0x%x ntp time:%v", addr, srcIA, t)
	return []uint8(t.Format("2006-01-02 15:04:05 -0700 MST")), 0
}
