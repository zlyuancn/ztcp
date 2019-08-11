/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/5/1
   Description :
-------------------------------------------------
*/

package server

import (
    "fmt"
    "github.com/zlyuancn/ztcp/client"
    "github.com/zlyuancn/ztcp/config"
    "net"
    "sync"
    "sync/atomic"
)

type Server struct {
    status config.ServerStatus
    opts   *Options
    mx     sync.Mutex
}

func NewServer(opts ...Option) (*Server, error) {
    options := newOptions(opts...)

    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", options.BindIP, options.BindPort))
    if err != nil {
        return nil, err
    }

    server := &Server{
        status: config.ServerListening,
        opts:   options,
    }

    options.Listener = listener
    options.Clients = make(clientStorage, options.InitClientCapacity)
    options.ClientConnectObserves = append([]client.ClientConnectObserve{func(c *client.Client) {
        server.addClient(c)
    }}, options.ClientConnectObserves...)
    options.ClientCloseObserves = append([]client.ClientCloseObserve{func(c *client.Client, err error) {
        server.removeClient(c)
    }}, options.ClientCloseObserves...)

    go func(m *Server) {
        for m.IsListening() {
            conn, err := listener.Accept()
            if err != nil {
                continue
            }
            go m.connectedHandler(conn)
        }
    }(server)
    return server, nil
}

func (m *Server) connectedHandler(conn net.Conn) {
    _, _ = client.NewClient(
        client.WithServerClient(conn),
        client.WithHeartbeatInterval(m.opts.HeartbeatCheckTime),
        client.WithClientConnectObserves(m.opts.ClientConnectObserves...),
        client.WithClientCloseObserves(m.opts.ClientCloseObserves...),
        client.WithClientSendDataObserves(m.opts.ClientSendDataObserves...),
        client.WithClientGetDataObserves(m.opts.ClientGetDataObserves...),
    )
}

func (m *Server) Options() *Options {
    return m.opts
}

func (m *Server) Addr() net.Addr {
    return m.opts.Listener.Addr()
}

func (m *Server) Status() config.ServerStatus {
    return config.ServerStatus(atomic.LoadInt32((*int32)(&m.status)))
}

func (m *Server) IsClosed() bool {
    return m.Status() == config.ServerClosed
}

func (m *Server) IsListening() bool {
    return m.Status() == config.ServerListening
}

func (m *Server) Close() error {
    err := m.opts.Listener.Close()
    if err == nil {
        atomic.StoreInt32((*int32)(&m.status), int32(config.ServerClosed))
    }
    return err
}

func (m *Server) CloseAllClient() (err error) {
    clients :=func() clientStorage{
        m.mx.Lock()
        defer m.mx.Unlock()
        return m.opts.Clients.Copy()
    }()

    for clientid, c := range clients {
        go func(clientId uint64, c *client.Client) {
            _ = c.Close()
        }(clientid, c)
    }
    return nil
}

func (m *Server) SendAll(data []byte) (err error) {
    clients :=func() clientStorage{
        m.mx.Lock()
        defer m.mx.Unlock()
        return m.opts.Clients.Copy()
    }()

    for clientid, c := range clients {
        go func(clientId uint64, c *client.Client) {
            _ = c.Send(data)
        }(clientid, c)
    }
    return nil
}

func (m *Server) addClient(c *client.Client) {
    go func() {
        m.mx.Lock()
        defer m.mx.Unlock()
        m.opts.Clients[c.GetId()] = c
    }()
}

func (m *Server) removeClient(c *client.Client) {
    go func() {
        m.mx.Lock()
        defer m.mx.Unlock()
        delete(m.opts.Clients, c.GetId())
    }()
}
