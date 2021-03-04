package main

import (
	"net"
	"time"
)

type libConnDeadlineRefresher struct {
	conn net.Conn
	idle time.Duration
}

func (me libConnDeadlineRefresher) Write(p []byte) (int, error) {
	me.conn.SetDeadline(time.Now().Add(me.idle))
	return len(p), nil
}
