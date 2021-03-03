package main

import (
	"log"

	"github.com/libp2p/go-libp2p-core/network"
)

func init() {
	handlers["None"] = func(
		cfg struct {
		},
		logger *log.Logger,
	) func(str network.Stream) {
		return func(str network.Stream) {
			str.Close()
		}
	}
}
