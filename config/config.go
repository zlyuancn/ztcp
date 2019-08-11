/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/5/2
   Description :
-------------------------------------------------
*/

package config

import (
    "time"
)

//数据传输缓存大小
var DefaultDataBuffSize = 1024 * 64
//一个包传输最大允许缓存大小
var DefaultDataPackageBuffSize = 1024 * 1024 * 64

// 心跳时间精度
var DefaultHeartbeatPrecision time.Duration = 1e9

// 默认信任消息
var DefaultTrustMsg = []byte("hello ztcp")
// 等待信任时间
var DefaultWaitTrustTime time.Duration = 5e9

const (
    //客户端Id占用字节数
    DataClientIdLength = 8
    //数据头占用字节数
    DataHeaderLength = 4
)
