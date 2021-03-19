// Copyright 2021-2021 The jdh99 Authors. All rights reserved.
// 配置文件
// Authors: jdh99 <jdh821@163.com>

package config

import "fmt"

// 系统参数
const (
	LocalIA   = 0x2141000000000004
	LocalIP   = "0.0.0.0"
	LocalPort = 12930
)

var LocalPwd string

func init() {
	fmt.Println("please input password:")
	_, err := fmt.Scanln(&LocalPwd)
	if err != nil {
		panic(err)
	}
	fmt.Println(LocalPwd)
}
