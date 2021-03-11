package forwardtcp

import (
	"carb/lib/deadline"
	"carb/lib/whitelist"
	"io"
	"log"
	"net"
	"path"
	"reflect"
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
	To          string
	IdleTimeout time.Duration
}

// Init ...
func (me *T) Init(hst host.Host, logger *log.Logger) {
	var filter whitelist.T
	filter.Add(me.WhiteList...)
	hst.SetStreamHandler(protocol.ID(
		path.Base(reflect.TypeOf(me).Elem().PkgPath()),
	), func(str network.Stream) {
		var conn net.Conn
		var err error
		var con deadline.T
		var wg sync.WaitGroup

		defer str.Close()

		if filter.Query(str.Conn().RemotePeer()) == false {
			logger.Println("[Banned]:", "peer is not in whitelist", "(id):", str.Conn().RemotePeer())
			return
		}

		if conn, err = net.DialTimeout("tcp", me.To, time.Second); err != nil {
			logger.Println("[Error]:", err)
			return
		}

		con.Conn = conn
		con.IdleTimeout = me.IdleTimeout
		con.Refresh()
		wg.Add(2)

		defer con.Close()
		defer wg.Wait()

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
	})
}
