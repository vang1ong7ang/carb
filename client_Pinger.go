package main

import (
	"context"
	"log"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
)

func init() {
	clients["Pinger"] = func(
		cfg struct {
			TimeInterval time.Duration
			Target       peer.ID
			PrintRTT     bool
		},
		logger *log.Logger,
	) func(hst host.Host) {
		return func(hst host.Host) {
			for {
				time.Sleep(cfg.TimeInterval)
				if result := <-ping.Ping(context.Background(), hst, cfg.Target); result.Error != nil {
					logger.Println("[Ping]:", result.Error)
				} else if cfg.PrintRTT {
					logger.Println("[RTT]:", cfg.Target, ":", result.RTT)
				}
			}
		}
	}
}
