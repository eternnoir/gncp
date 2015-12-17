package gncp

import (
	"net"
	"sync"
)

type ConnPool interface {
	Get() (net.Conn, error)
	Close() error
}

type GncpPool struct {
	lock       sync.Mutex
	conns      chan net.Conn
	minConnNum int
	maxConnNum int

	connCreater func() (net.Conn, error)
}

func NewPool(minConn, maxConn int, connCreater func() (net.Conn, error)) (*GncpPool, error) {
	pool := &GncpPool{}
	pool.minConnNum = minConn
	pool.maxConnNum = maxConn
	pool.connCreater = connCreater
	err := pool.init()
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func (p *GncpPool) init() error {
	return nil
}

func (p *GncpPool) createConn() (net.Conn, error) {
	return nil, nil
}
