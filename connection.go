package gncp

import (
	"net"
)

type CpConn struct {
	net.Conn
	pool *GncpPool
}
