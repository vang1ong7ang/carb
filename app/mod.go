package app

import (
	"carb/ext"
	"carb/lib/address"
	"carb/lib/sk"
	"context"
	"encoding/json"
	"log"

	"github.com/libp2p/go-libp2p"
	relay "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	quic "github.com/libp2p/go-libp2p-quic-transport"
)

// T ...
type T struct {
	SK        *sk.T            `json:"sk"`
	LISTEN    []string         `json:"listen"`
	RELAYMODE []relay.RelayOpt `json:"relaymode"`
	PEER      []struct {
		ID      peer.ID    `json:"id"`
		ADDRESS *address.T `json:"address"`
	} `json:"peer"`
	EXT []struct {
		ENABLE   bool            `json:"enable"`
		LOGLEVEL int             `json:"loglevel"`
		ID       string          `json:"id"`
		CONFIG   json.RawMessage `json:"config"`
	} `json:"ext"`
}

// Start ...
func (me *T) Start() {
	hst, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(me.LISTEN...),
		libp2p.Identity(me.SK),
		libp2p.EnableRelay(me.RELAYMODE...),
		libp2p.NATPortMap(),
		libp2p.Transport(quic.NewTransport),
		libp2p.DefaultTransports,
		libp2p.Routing(func(hst host.Host) (routing.PeerRouting, error) {
			return dht.New(context.Background(), hst)
		}),
		libp2p.EnableAutoRelay(),
	)

	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range me.PEER {
		hst.Peerstore().AddAddr(v.ID, v.ADDRESS, peerstore.AddressTTL)
	}

	hdlmgt := new(ext.T)
	hdlmgt.Init()

	for _, v := range me.EXT {
		if v.ENABLE == false {
			continue
		}
		go hdlmgt.Get(v.ID, v.CONFIG, hst, v.LOGLEVEL)
	}
}
