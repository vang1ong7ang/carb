package listentcp

import (
	"carb/lib/deadline"
	"context"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

// T ...
type T struct {
	Target      peer.ID
	Listen      string
	IdleTimeout time.Duration
}

// Init ...
func (me *T) Init(hst host.Host, logger *log.Logger) {
	var listener net.Listener
	var err error

	if listener, err = net.Listen("tcp", me.Listen); err != nil {
		logger.Println("[Listen]:", err)
		return
	}

	defer listener.Close()

	for {
		func() {
			var conn net.Conn
			var str network.Stream
			var con deadline.T
			var wg sync.WaitGroup

			if conn, err = listener.Accept(); err != nil {
				logger.Println("[Accept]:", err)
				return
			}

			defer conn.Close()

			if str, err = hst.NewStream(context.Background(), me.Target, "forwardtcp"); err != nil {
				logger.Println("[NewStream]:", err)
				return
			}

			defer str.Close()
			defer wg.Wait()

			con.Conn = conn
			con.IdleTimeout = me.IdleTimeout
			con.Refresh()
			wg.Add(2)

			go func() {
				var n int64
				var err error

				defer wg.Done()
				defer str.CloseRead()

				n, err = io.Copy(&con, str)
				logger.Println("[Str -> Con]:", n, "(err):", err)
			}()

			go func() {
				var n int64
				var err error

				defer wg.Done()
				defer str.CloseWrite()

				n, err = io.Copy(str, &con)
				logger.Println("[Con -> Str]:", n, "(err):", err)
			}()

		}()
	}
}
