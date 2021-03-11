package peerlog

import (
	"log"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
)

// T ...
type T struct {
	TimeInterval time.Duration
}

// Init ...
func (me *T) Init(hst host.Host, logger *log.Logger) {
	var text []byte
	var err error
	for {
		for _, v := range hst.Peerstore().Peers() {
			if text, err = v.MarshalText(); err != nil {
				logger.Println("[Error]:", err)
				continue
			}
			logger.Println("[Peer]:", string(text))
		}
	}
}
