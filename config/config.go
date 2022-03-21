// Copyright 2021-2021 The jdh99 Authors. All rights reserved.
// 配置文件
// Authors: jdh99 <jdh821@163.com>

package config

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"strconv"
)

// 配置文件解析后的结构
type tConfigJson struct {
	LocalIA   string `json:"LocalIA"`
	LocalIP   string `json:"LocalIP"`
	LocalPort int    `json:"LocalPort"`
	CoreIA    string `json:"CoreIA"`
	CoreIP    string `json:"CoreIP"`
	CorePort  int    `json:"CorePort"`
}

// LocalIA 本机IA地址
var LocalIA uint32

// LocalIP 本机IP
var LocalIP uint32
var LocalIPStr string

// LocalPort 本机端口
var LocalPort uint16

// CoreIA 核心网IA地址
var CoreIA uint32

// CoreIP 核心网IP
var CoreIP uint32
var CoreIPStr string

// CorePort 核心网端口
var CorePort uint16

func init() {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	configJson := &tConfigJson{}
	err = json.Unmarshal(data, &configJson)
	if err != nil {
		panic(err)
	}

	temp, err := strconv.ParseUint(configJson.LocalIA, 0, 32)
	if err != nil {
		panic(err)
	}
	LocalIA = uint32(temp)

	LocalPort = uint16(configJson.LocalPort)
	LocalIPStr = configJson.LocalIP
	addr := net.UDPAddr{IP: net.ParseIP(LocalIPStr), Port: int(LocalPort)}
	ipValue := addr.IP.To4()
	LocalIP = (uint32(ipValue[0]) << 24) + (uint32(ipValue[1]) << 16) + (uint32(ipValue[2]) << 8) + uint32(ipValue[3])

	temp, err = strconv.ParseUint(configJson.CoreIA, 0, 32)
	if err != nil {
		panic(err)
	}
	CoreIA = uint32(temp)

	CorePort = uint16(configJson.CorePort)
	CoreIPStr = configJson.CoreIP
	addr = net.UDPAddr{IP: net.ParseIP(CoreIPStr), Port: int(CorePort)}
	ipValue = addr.IP.To4()
	CoreIP = (uint32(ipValue[0]) << 24) + (uint32(ipValue[1]) << 16) + (uint32(ipValue[2]) << 8) + uint32(ipValue[3])
}
