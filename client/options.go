/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/8/10
   Description :
-------------------------------------------------
*/

package client

import (
    "context"
    "net"
    "time"
)

//心跳发送时间(推荐为心跳检测时间的 2/5
var DefaultHeartbeatInterval time.Duration = 16e9

type Option func(opts *Options)

type ClientConnectObserve func(c *Client)
type ClientCloseObserve func(c *Client, err error)
type ClientSendDataObserve func(c *Client, data []byte)
type ClientGetDataObserve func(c *Client, data []byte)

type Options struct {
    IsServerClient bool
    Conn           net.Conn
    // bind本地端口
    BindPort int
    // 要连接的地址
    ConnectAddr string
    // 可以用于连接超时
    ConnectContext context.Context
    // 客户端连接观察者
    ClientConnectObserves []ClientConnectObserve
    // 客户端关闭观察者
    ClientCloseObserves []ClientCloseObserve
    // 发送数据观察者
    ClientSendDataObserves []ClientSendDataObserve
    // 获取数据观察者
    ClientGetDataObserves []ClientGetDataObserve
    // 心跳间隔时间
    HeartbeatInterval time.Duration
}

func newOptions(opts ...Option) *Options {
    opt := &Options{
        HeartbeatInterval: DefaultHeartbeatInterval,
        ConnectContext:    context.Background(),
    }

    for _, o := range opts {
        o(opt)
    }
    return opt
}

func WithServerClient(conn net.Conn) Option {
    return func(opts *Options) {
        opts.IsServerClient = true
        opts.Conn = conn
    }
}

func WithBindPort(port int) Option {
    return func(opts *Options) {
        opts.BindPort = port
    }
}

func WithConnectAddr(addr string) Option {
    return func(opts *Options) {
        opts.ConnectAddr = addr
    }
}

func WithConnectContext(ctx context.Context) Option {
    return func(opts *Options) {
        opts.ConnectContext = ctx
    }
}

func WithClientConnectObserves(observers ...ClientConnectObserve) Option {
    return func(opts *Options) {
        opts.ClientConnectObserves = append(opts.ClientConnectObserves, observers...)
    }
}

func WithClientCloseObserves(observers ...ClientCloseObserve) Option {
    return func(opts *Options) {
        opts.ClientCloseObserves = append(opts.ClientCloseObserves, observers...)
    }
}

func WithClientSendDataObserves(observers ...ClientSendDataObserve) Option {
    return func(opts *Options) {
        opts.ClientSendDataObserves = append(opts.ClientSendDataObserves, observers...)
    }
}

func WithClientGetDataObserves(observers ...ClientGetDataObserve) Option {
    return func(opts *Options) {
        opts.ClientGetDataObserves = append(opts.ClientGetDataObserves, observers...)
    }
}

func WithHeartbeatInterval(interval time.Duration) Option {
    return func(opts *Options) {
        opts.HeartbeatInterval = interval
    }
}
