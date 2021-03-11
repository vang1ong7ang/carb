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
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	quic "github.com/libp2p/go-libp2p-quic-transport"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	if hst, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(cfg.LISTEN...),
		libp2p.Identity(cfg.sk()),
		libp2p.EnableRelay(cfg.RELAYMODE...),
		libp2p.NATPortMap(),
		libp2p.Transport(quic.NewTransport),
		libp2p.DefaultTransports,
		libp2p.Routing(func(hst host.Host) (routing.PeerRouting, error) {
			return dht.New(context.Background(), hst)
		}),
		libp2p.EnableAutoRelay(),
	); err != nil {
		log.Fatalln("[P2P]:", err)
	} else {
		log.Println("[HostID]:", hst.ID())

		for k, v := range cfg.PEERS {
			// for _, vv := range v {
			// 	// var ma multiaddr.Multiaddr
			// 	ma, _ := multiaddr.NewMultiaddr(vv)
			// 	st, _ := json.Marshal(ma)
			// 	log.Println("CCCC", string(st))
			// 	log.Println("DDDD", st)
			// 	log.Println("EEEE", vv)
			// 	ss := string(st)
			// 	_ = ss
			// 	log.Println(json.Unmarshal([]byte(vv), &ma))

			// }
			if addr, err := multiaddr.NewMultiaddr(v.Address); err != nil {
				log.Fatalln("[Address]:", err)
			} else {
				hst.Peerstore().AddAddr(v.ID, addr, peerstore.TempAddrTTL)
			}
			log.Println("[Peer]:", k, ":", v)
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
