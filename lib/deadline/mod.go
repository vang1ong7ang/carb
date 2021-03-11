package deadline

import (
	"net"
	"time"
)

// T ...
type T struct {
	net.Conn
	IdleTimeout time.Duration
}

// Write ...
func (me *T) Write(p []byte) (int, error) {
	me.Refresh()
	return me.Conn.Write(p)
}

// Read ...
func (me *T) Read(p []byte) (int, error) {
	me.Refresh()
	return me.Conn.Read(p)
}

// Refresh ...
func (me *T) Refresh() {
	me.SetDeadline(time.Now().Add(me.IdleTimeout))
}
