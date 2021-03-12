package peerlog

import (
	"carb/lib/logger"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
)

// T ...
type T struct {
	TimeInterval time.Duration
}

// Init ...
func (me *T) Init(hst host.Host, log *logger.T) {
	for {
		for _, v := range hst.Peerstore().Peers() {
			text, err := v.MarshalText()
			if err != nil {
				log.Log(1, "fail", map[string]interface{}{"err": err})
				continue
			}
			log.Log(2, "peer", map[string]interface{}{"peer": string(text)})
		}
	}
}
