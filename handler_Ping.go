package main

import (
	"log"

	"github.com/libp2p/go-libp2p-core/network"
)

func init() {
	handlers["Ping"] = func(
		cfg struct {
			PrintLog bool
		},
		logger *log.Logger,
	) func(str network.Stream) {
		return func(str network.Stream) {
			if cfg.PrintLog {
				logger.Println("[RECVFROM]:", str.Conn().RemotePeer())
			}
			str.Write([]byte("pong"))
			str.Close()
		}
	}
}
