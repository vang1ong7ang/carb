package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	if hst, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(cfg.LISTEN),
		libp2p.Identity(cfg.sk()),
		libp2p.EnableRelay(cfg.RELAYMODE...),
	); err != nil {
		log.Fatalln("[P2P]:", err)
	} else {
		log.Println("[HostID]:", hst.ID())

		for k, v := range cfg.PEERS {
			if id, err := peer.IDB58Decode(k); err != nil {
				continue
			} else if addr, err := multiaddr.NewMultiaddr(v); err != nil {
				continue
			} else {
				hst.Peerstore().AddAddr(id, addr, peerstore.PermanentAddrTTL)
				log.Println("[Peer]:", id)
			}
		}

		for k, v := range cfg.PROTOCOLS {
			cfg := reflect.New(reflect.TypeOf(handlers[v.ID]).In(0))
			if err := json.Unmarshal(v.CONFIG, cfg.Interface()); err != nil {
				log.Fatalln(err)
			}
			hst.SetStreamHandler(
				v.ID,
				reflect.ValueOf(handlers[v.ID]).
					Call([]reflect.Value{
						cfg.Elem(),
						reflect.ValueOf(log.New(os.Stderr, fmt.Sprintf("{{%s}}: ", v.ID), log.LstdFlags|log.Lmsgprefix)),
					})[0].Interface().(func(network.Stream)),
			)
			log.Println("[PROTOCOL]:", k, "*", v.ID)
		}

		for k, v := range cfg.CLIENTS {
			cfg := reflect.New(reflect.TypeOf(clients[v.ID]).In(0))
			if err := json.Unmarshal(v.CONFIG, cfg.Interface()); err != nil {
				log.Fatalln(err)
			}
			go reflect.ValueOf(clients[v.ID]).
				Call([]reflect.Value{
					cfg.Elem(),
					reflect.ValueOf(log.New(os.Stderr, fmt.Sprintf("{{%s}}: ", v.ID), log.LstdFlags|log.Lmsgprefix)),
				})[0].Interface().(func(host.Host))(hst)
			log.Println("[CLIENT]:", k, "*", v.ID)
		}
	}

	select {}
}
