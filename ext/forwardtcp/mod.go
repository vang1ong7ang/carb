package forwardtcp

import (
	"carb/lib/deadline"
	"carb/lib/logger"
	"carb/lib/whitelist"
	"io"
	"net"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

// T ...
type T struct {
	WhiteList   []peer.ID
	ProtocolID  protocol.ID
	To          string
	IdleTimeout time.Duration
}

// Init ...
func (me *T) Init(hst host.Host, log *logger.T) {
	filter := new(whitelist.T)
	filter.Add(me.WhiteList...)

	hst.SetStreamHandler(me.ProtocolID, func(str network.Stream) {

		defer str.Close()

		if filter.Query(str.Conn().RemotePeer()) == false {
			log.Log(2, "ban", map[string]interface{}{"id": str.Conn().RemotePeer()})
			return
		}

		conn, err := net.DialTimeout("tcp", me.To, time.Second)

		if err != nil {
			log.Log(3, "dial", map[string]interface{}{"err": err})
			return
		}

		defer conn.Close()

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
	})
}
