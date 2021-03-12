package listentcp

import (
	"carb/lib/deadline"
	"carb/lib/logger"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

// T ...
type T struct {
	Target      peer.ID
	ProtocolID  protocol.ID
	Listen      string
	IdleTimeout time.Duration
}

// Init ...
func (me *T) Init(hst host.Host, log *logger.T) {
	listener, err := net.Listen("tcp", me.Listen)
	if err != nil {
		log.Log(1, "listen", map[string]interface{}{"err": err})
		return
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Log(2, "accept", map[string]interface{}{"err": err})
			continue
		}

		go func() {
			defer conn.Close()

			str, err := hst.NewStream(context.Background(), me.Target, me.ProtocolID)

			if err != nil {
				log.Log(3, "new stream", map[string]interface{}{"err": err})
				return
			}

			defer str.Close()

			con := &deadline.T{
				Conn:        conn,
				IdleTimeout: me.IdleTimeout,
			}

			con.Refresh()

			wg := new(sync.WaitGroup)
			wg.Add(2)

			defer wg.Wait()

			go func() {
				defer wg.Done()
				defer str.CloseRead()

				n, err := io.Copy(con, str)
				log.Log(4, "str->con", map[string]interface{}{"n": n, "err": err})
			}()

			go func() {
				defer wg.Done()
				defer str.CloseWrite()

				n, err := io.Copy(str, con)
				log.Log(4, "con->str", map[string]interface{}{"n": n, "err": err})
			}()
		}()
	}
}
