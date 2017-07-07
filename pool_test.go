package gncp

import (
	"fmt"
	"testing"
	"time"

	"context"

	"github.com/stretchr/testify/assert"

	"net"
)

var (
	Host = "127.0.0.1"
	Port = "5566"
)

func init() {
	go exampleServer()
	time.Sleep(1 * time.Second)
}

func TestCreateConnpool(t *testing.T) {
	assert := assert.New(t)
	pool, err := NewPool(1, 10, connCreator)
	if err != nil {
		assert.Fail("Init conn pool fail", err)
	}
	assert.NotEmpty(pool)
}

func TestNewPoolConnNumError(t *testing.T) {
	assert := assert.New(t)
	_, err := NewPool(-1, 10, connCreator)
	if err != nil {
		return
	}
	assert.Fail("Should be fail.")
}

func TestNewPoolConnNumError2(t *testing.T) {
	assert := assert.New(t)
	_, err := NewPool(2, 1, connCreator)
	if err != nil {
		return
	}
	assert.Fail("Should be fail.")
}

func TestGetConn(t *testing.T) {
	assert := assert.New(t)
	pool, err := NewPool(1, 3, connCreator)
	if err != nil {
		assert.Fail("Init conn pool fail.")
		return
	}
	conn1, err := pool.Get()
	if err != nil {
		assert.Fail("Get conn fail.")
	}
	_, err = conn1.Write([]byte("Test conn1"))
	if err != nil {
		assert.Fail("Write message fail.")
	}
}

func TestGetConnTimeout(t *testing.T) {
	assert := assert.New(t)
	pool, err := NewPool(1, 3, connCreator)
	if err != nil {
		assert.Fail("Init conn pool fail.")
		return
	}

	for i := 0; i < 3; i++ {
		conn, err := pool.Get()
		if err != nil {
			assert.Fail("Get conn fail.")
		}
		_, err = conn.Write([]byte("Test conn1"))
		if err != nil {
			assert.Fail("Write message fail.")
		}
	}

	_, err = pool.GetWithTimeout(time.Duration(1) * time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Fail("Need get timeout error.")
}

func TestGetConnContxt(t *testing.T) {
	assert := assert.New(t)
	pool, err := NewPool(1, 3, connCreator)
	if err != nil {
		assert.Fail("Init conn pool fail.")
		return
	}

	// Get all connection from pool
	for i := 0; i < 3; i++ {
		conn, err := pool.Get()
		if err != nil {
			assert.Fail("Get conn fail.")
		}
		_, err = conn.Write([]byte("Test conn1"))
		if err != nil {
			assert.Fail("Write message fail.")
		}
	}
	ctx := context.Background()
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(1)*time.Second)
	defer cancel()
	_, err = pool.GetWithContext(ctxTimeout)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Fail("Need get timeout error.")
}

func TestRemoveConn(t *testing.T) {
	assert := assert.New(t)
	pool, err := NewPool(1, 3, connCreator)
	if err != nil {
		assert.Fail("Init conn pool fail.")
	}
	conn1, err := pool.Get()
	err = pool.Remove(conn1)
	if err != nil {
		assert.Fail("Cannot remoce connection.")
	}
	err = conn1.Close()
	if err != nil {
		return
	}
	assert.Fail("Need connection already removed error.")
}
func connCreator() (net.Conn, error) {
	return net.Dial("tcp", Host+":"+Port)
}

func exampleServer() {
	l, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go func() {
			buffer := make([]byte, 1024)
			conn.Read(buffer)
		}()
	}
}
