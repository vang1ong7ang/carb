package pinger

import (
	"carb/lib/logger"
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
)

// T ...
type T struct {
	Target       peer.ID
	TimeInterval time.Duration
}

// Init ...
func (me *T) Init(hst host.Host, log *logger.T) {
	for {
		time.Sleep(me.TimeInterval)
		result := <-ping.Ping(context.Background(), hst, me.Target)
		if result.Error != nil {
			log.Log(1, "fail", map[string]interface{}{"err": result.Error})
			return
		}
		log.Log(2, "ping", map[string]interface{}{"tgt": me.Target, "rtt": result.RTT})
	}
}
