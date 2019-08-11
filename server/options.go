/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/8/10
   Description :
-------------------------------------------------
*/

package server

import (
    "github.com/zlyuancn/ztcp/client"
    "net"
    "time"
)

//初始客户端容量
var DefaultInitClientCapacity = 1000
//心跳检测时间
var DefaultHeartbeatCheckTime time.Duration = 40e9

// 保存所有已连接成功的客户端
type clientStorage map[uint64]*client.Client

// 拷贝一个副本
func (m clientStorage) Copy() clientStorage {
    a := make(clientStorage, len(m))
    for k, v := range (m) {
        a[k] = v
    }
    return a
}

type Option func(opts *Options)

type Options struct {
    Listener net.Listener
    Clients  clientStorage

    BindIP   string
    BindPort int
    // 初始客户端容量
    InitClientCapacity int
    // 客户端连接观察者
    ClientConnectObserves []client.ClientConnectObserve
    // 客户端关闭观察者
    ClientCloseObserves []client.ClientCloseObserve
    // 发送数据观察者
    ClientSendDataObserves []client.ClientSendDataObserve
    // 获取数据观察者
    ClientGetDataObserves []client.ClientGetDataObserve
    // 检查心跳时间
    HeartbeatCheckTime time.Duration
}

func newOptions(opts ...Option) *Options {
    opt := &Options{
        InitClientCapacity: DefaultInitClientCapacity,
        HeartbeatCheckTime: DefaultHeartbeatCheckTime,
    }

    for _, o := range opts {
        o(opt)
    }
    return opt
}

func WithBindIP(bindip string) Option {
    return func(opts *Options) {
        opts.BindIP = bindip
    }
}

func WithBindPort(port int) Option {
    return func(opts *Options) {
        opts.BindPort = port
    }
}

func WithClientCapacity(capacity int) Option {
    return func(opts *Options) {
        opts.InitClientCapacity = capacity
    }
}

func WithClientConnectObserves(observers ...client.ClientConnectObserve) Option {
    return func(opts *Options) {
        opts.ClientConnectObserves = append(opts.ClientConnectObserves, observers...)
    }
}

func WithClientCloseObserves(observers ...client.ClientCloseObserve) Option {
    return func(opts *Options) {
        opts.ClientCloseObserves = append(opts.ClientCloseObserves, observers...)
    }
}

func WithClientSendDataObserves(observers ...client.ClientSendDataObserve) Option {
    return func(opts *Options) {
        opts.ClientSendDataObserves = append(opts.ClientSendDataObserves, observers...)
    }
}

func WithClientGetDataObserves(observers ...client.ClientGetDataObserve) Option {
    return func(opts *Options) {
        opts.ClientGetDataObserves = append(opts.ClientGetDataObserves, observers...)
    }
}

func WithHeartbeatCheckTime(checktime time.Duration) Option {
    return func(opts *Options) {
        opts.HeartbeatCheckTime = checktime
    }
}
