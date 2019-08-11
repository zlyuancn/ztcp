/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/8/10
   Description :
-------------------------------------------------
*/

package config

type (
    ServerStatus int32
    ClientStatus int32
)

const (
    //服务端已关闭
    ServerClosed ServerStatus = iota
    //服务端监听中
    ServerListening
)

const (
    //客户端已关闭
    ClientClosed ClientStatus = iota
    //客户端连接中
    ClientConnecting
    //等待信任
    ClientWaitTrust
    //连接成功
    ClientConnected
)
