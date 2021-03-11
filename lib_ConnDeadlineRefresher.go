package main

import (
	"net"
	"time"
)

type libConnDeadlineRefresher struct {
	net.Conn
	idle time.Duration
}

func (me *libConnDeadlineRefresher) Write(p []byte) (int, error) {
	me.Refresh()
	return me.Conn.Write(p)
}

func (me *libConnDeadlineRefresher) Read(p []byte) (int, error) {
	me.Refresh()
	return me.Conn.Read(p)
}

func (me *libConnDeadlineRefresher) Refresh() {
	me.SetDeadline(time.Now().Add(me.idle))
}
