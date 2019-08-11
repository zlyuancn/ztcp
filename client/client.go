/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/5/1
   Description :
-------------------------------------------------
*/

package client

import (
    "bytes"
    "errors"
    "github.com/zlyuancn/zassert"
    "github.com/zlyuancn/ztcp/config"
    "github.com/zlyuancn/ztcp/utils"
    "net"
    "sync"
    "sync/atomic"
    "time"
)

type Client struct {
    clientId      uint64
    status        config.ClientStatus
    opts          *Options
    mx            sync.Mutex
    heartbeatTime *utils.HeartbeatTime
}

func NewClient(opts ...Option) (*Client, error) {
    options := newOptions(opts...)

    c := &Client{
        opts:   options,
        status: config.ClientConnecting,
    }

    if options.IsServerClient {
        c.connectedHandler()
        return c, nil
    }

    if options.ConnectAddr == "" {
        return nil, zassert.AssertError{"未设置 ConnectAddr ,请使用 WithConnectAddr(addr)"}
    }

    go func(m *Client) {
        var d net.Dialer
        d.LocalAddr = &net.TCPAddr{Port: m.opts.BindPort}

        conn, err := d.DialContext(m.opts.ConnectContext, "tcp", m.opts.ConnectAddr)
        if err != nil {
            m.closedHandler(err)
            return
        }
        m.opts.Conn = conn
        m.connectedHandler()
    }(c)

    return c, nil
}

func (m *Client) Options() *Options {
    return m.opts
}

func (m *Client) RemoteAddr() net.Addr {
    return m.opts.Conn.RemoteAddr()
}

func (m *Client) LocalAddr() net.Addr {
    return m.opts.Conn.LocalAddr()
}

func (m *Client) GetId() uint64 {
    return m.clientId
}

func (m *Client) Status() config.ClientStatus {
    return config.ClientStatus(atomic.LoadInt32((*int32)(&m.status)))
}

func (m *Client) IsClosed() bool {
    return m.Status() == config.ClientClosed
}

func (m *Client) IsConnected() bool {
    return m.Status() == config.ClientConnected
}

func (m *Client) Close() error {
    if m.opts.Conn == nil {
        return nil
    }
    return m.opts.Conn.Close()
}

func (m *Client) Send(data []byte) (err error) {
    if m.Status() != config.ClientConnected {
        return zassert.AssertError{"Client 非 ClientConnected 状态时不能使用 Send"}
    }
    if len(data) == 0 {
        return
    }

    m.mx.Lock()
    defer m.mx.Unlock()

    m.notifyClientSendData(m, data)

    dataHeader := utils.Uint32ToBytes(uint32(len(data)))
    if _, err = m.opts.Conn.Write(dataHeader); err != nil {
        return nil
    }
    m.heartbeatTime.RefHeartbeat()

    _, err = m.opts.Conn.Write(data)
    return err
}

func (m *Client) waitTrust(trust chan struct{}, distrust chan error) {
    sendTrustMsg := func() error {
        _, err := m.opts.Conn.Write(config.DefaultTrustMsg)
        return err
    }
    waitTrustMsg := func() error {
        buff, err := utils.WaitConnData(m.opts.Conn, len(config.DefaultTrustMsg))
        if err != nil {
            return err
        }
        if !bytes.Equal(buff, config.DefaultTrustMsg) {
            return errors.New("信任消息错误")
        }
        return nil
    }

    if m.opts.IsServerClient {
        if err := waitTrustMsg(); err != nil {
            distrust <- err
            return
        }

        if err := sendTrustMsg(); err != nil {
            distrust <- err
            return
        }

        m.clientId = utils.AutoClientID.Next()
        _, err := m.opts.Conn.Write(utils.Uint64ToBytes(m.clientId))
        if err != nil {
            distrust <- err
            return
        }

    } else {
        if err := sendTrustMsg(); err != nil {
            distrust <- err
            return
        }
        if err := waitTrustMsg(); err != nil {
            distrust <- err
            return
        }

        buff, err := utils.WaitConnData(m.opts.Conn, config.DataClientIdLength)
        if err != nil {
            distrust <- err
            return
        }
        m.clientId = utils.BytesToUint64(buff)
    }

    trust <- struct{}{}
}

func (m *Client) connectedHandler() {
    m.changeStatus(config.ClientWaitTrust)

    var trust = make(chan struct{}, 1)
    var distrust = make(chan error, 1)
    var trust_time = time.NewTicker(config.DefaultWaitTrustTime)

    go m.waitTrust(trust, distrust)

    select {
    case <-trust:
        trust_time.Stop()
    case err := <-distrust:
        m.closedHandler(err)
        return
    case <-trust_time.C:
        m.closedHandler(errors.New("超过最大信任等待时间"))
        return
    }

    m.heartbeatTime = utils.NewHeartbeatTime(m.opts.HeartbeatInterval, config.DefaultHeartbeatPrecision, m.heartbeatFunc)
    m.changeStatus(config.ClientConnected)
    m.notifyClientConnect(m)
    m.received()
}

func (m *Client) received() {
    for m.Status() == config.ClientConnected {
        data, err := utils.WaitConnFullData(m.opts.Conn)
        if err != nil {
            m.closedHandler(nil)
            return
        }

        m.heartbeatTime.RefHeartbeat()
        if data != nil && len(data) > 0 {
            m.notifyClientGetData(m, data)
        }
    }
}

func (m *Client) heartbeatFunc(timer *utils.HeartbeatTime) {
    if m.opts.IsServerClient {
        _ = m.Close()
    } else {
        m.mx.Lock()
        defer m.mx.Unlock()

        // 发送一个空数据来表示心跳
        _, _ = m.opts.Conn.Write(make([]byte, config.DataHeaderLength))
    }
}

func (m *Client) closedHandler(err error) {
    if m.Status() != config.ClientClosed {
        m.changeStatus(config.ClientClosed)
        if m.heartbeatTime != nil {
            m.heartbeatTime.Stop()
        }

        m.notifyClientClose(m, err)
    }
}

func (m *Client) changeStatus(status config.ClientStatus) {
    atomic.StoreInt32((*int32)(&m.status), int32(status))
}

func (m *Client) notifyClientConnect(c *Client) {
    for _, fn := range c.opts.ClientConnectObserves {
        fn(c)
    }
}
func (m *Client) notifyClientClose(c *Client, err error) {
    for _, fn := range c.opts.ClientCloseObserves {
        fn(c, err)
    }
}
func (m *Client) notifyClientSendData(c *Client, data []byte) {
    for _, fn := range c.opts.ClientSendDataObserves {
        fn(c, data)
    }
}
func (m *Client) notifyClientGetData(c *Client, data []byte) {
    for _, fn := range c.opts.ClientGetDataObserves {
        fn(c, data)
    }
}
