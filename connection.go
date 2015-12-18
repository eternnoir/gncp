package gncp

import (
	"net"
)

type CpConn struct {
	net.Conn
	pool *GncpPool
}

// Close will push connection back to connection pool. It will not close the real connection.
func (conn *CpConn) Close() error {
	return conn.pool.Put(conn.Conn)
}
