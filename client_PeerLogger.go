package main

import (
	"log"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
)

func init() {
	clients["PeerLogger"] = func(
		cfg struct {
			TimeInterval time.Duration
		},
		logger *log.Logger,
	) func(hst host.Host) {
		return func(hst host.Host) {
			for {
				time.Sleep(cfg.TimeInterval)
				logger.Println("[PeersStart]:")
				for k, v := range hst.Peerstore().Peers() {
					logger.Println(k, v)
				}
				logger.Println("[PeersEnd]:")
			}
		}
	}
}
