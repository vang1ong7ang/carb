package main

import (
	"context"
	"log"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

func init() {
	clients["HeartBeater"] = func(
		cfg struct {
			TimeInterval time.Duration
			Target       string
			PrintLog     bool
		},
		logger *log.Logger,
	) func(hst host.Host) {
		return func(hst host.Host) {
			for {
				time.Sleep(cfg.TimeInterval)
				if id, err := peer.IDB58Decode(cfg.Target); err != nil {
					logger.Println("[IDDecode]:", err)
				} else if str, err := hst.NewStream(context.Background(), id, "Ping"); err != nil {
					logger.Println("[NewStream]:", err)
				} else if n, err := str.Write([]byte("ping")); err != nil {
					log.Println("[WriteStream]", err, "{<n>}:", n)
					str.Close()
				} else {
					if cfg.PrintLog {
						logger.Println("[PING]:", id)
					}
					str.Close()
				}
			}
		}
	}
}
