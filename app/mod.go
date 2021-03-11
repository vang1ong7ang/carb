package app

import (
	"carb/hdl"
	"carb/lib/address"
	"carb/lib/sk"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/libp2p/go-libp2p"
	relay "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	quic "github.com/libp2p/go-libp2p-quic-transport"
)

// T ...
type T struct {
	SK        sk.T
	LISTEN    []string
	RELAYMODE []relay.RelayOpt
	PEERS     []struct {
		ID      peer.ID
		ADDRESS address.T
	}
	Handlers []struct {
		ENABLE bool
		LOG    bool
		ID     protocol.ID
		CONFIG json.RawMessage
	}
}

// Start ...
func (me *T) Start() {
	var hst host.Host
	var hdlmgt hdl.T
	var err error
	if hst, err = libp2p.New(
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
	); err != nil {
		return
	}

	for _, v := range me.PEERS {
		hst.Peerstore().AddAddr(v.ID, &v.ADDRESS, peerstore.AddressTTL)
	}

	hdlmgt.Init()

	for _, v := range me.Handlers {
		if v.ENABLE == false {
			continue
		}
		if v.LOG {
			go hdlmgt.Get(v.ID, v.CONFIG, hst, os.Stderr)
		} else {
			go hdlmgt.Get(v.ID, v.CONFIG, hst, ioutil.Discard)
		}
	}
	return
}
