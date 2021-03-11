package pinger

import (
	"context"
	"log"
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
func (me *T) Init(hst host.Host, logger *log.Logger) {
	var result ping.Result
	for {
		time.Sleep(me.TimeInterval)
		if result = <-ping.Ping(context.Background(), hst, me.Target); result.Error != nil {
			logger.Println("[Error]:", result.Error)
			return
		}
		logger.Println("[Report]:", "(Target):", me.Target, "(RTT):", result.RTT)
	}
}
