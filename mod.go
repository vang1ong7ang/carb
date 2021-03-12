package main

import (
	"carb/app"
	"carb/ext"
	"carb/lib/address"
	"carb/lib/sk"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	relay "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

// T ...
type T struct{}

// Addext ...
func (me *T) Addext(args struct {
	Config string `default:".carb.json" help:"carb config file path"`
	ID     string `default:"pinger" help:"handler protocol id"`
}) {
	app := me.open(args.Config)
	hdlmgt := new(ext.T)
	hdlmgt.Init()

	var handler struct {
		ENABLE   bool            `json:"enable"`
		LOGLEVEL int             `json:"loglevel"`
		ID       string          `json:"id"`
		CONFIG   json.RawMessage `json:"config"`
	}

	cfg, err := hdlmgt.GetConfig(args.ID)
	if err != nil {
		log.Fatalln(err)
	}
	handler.CONFIG = cfg
	handler.ENABLE = true
	handler.ID = args.ID
	app.EXT = append(app.EXT, handler)
	me.save(args.Config, app)
}

// Addlis ...
func (me *T) Addlis(args struct {
	Config  string `default:".carb.json" help:"carb config file path"`
	Address string `default:"/ip4/0.0.0.0/tcp/5517" help:"listen address in multiaddr format"`
}) {
	app := me.open(args.Config)
	app.LISTEN = append(app.LISTEN, args.Address)
	me.save(args.Config, app)
}

// Addpeer ...
func (me *T) Addpeer(args struct {
	Config  string `default:".carb.json" help:"carb config file path"`
	ID      string `default:"<EMPTY PEER ID>" help:"remote peer id"`
	Address string `default:"<EMPTY ADDRESS>" help:"peer address in multiaddr format"`
}) {
	app := me.open(args.Config)
	addr, err := multiaddr.NewMultiaddr(args.Address)
	if err != nil {
		log.Fatalln(err)
	}
	id, err := peer.IDB58Decode(args.ID)
	if err != nil {
		log.Fatalln(err)
	}
	app.PEER = append(app.PEER, struct {
		ID      peer.ID    "json:\"id\""
		ADDRESS *address.T "json:\"address\""
	}{
		ID:      id,
		ADDRESS: &address.T{Multiaddr: addr},
	})
	me.save(args.Config, app)
}

// Enablerelay ...
func (me *T) Enablerelay(args struct {
	Config  string `default:".carb.json" help:"carb config file path"`
	Address string `default:"/ip4/0.0.0.0/tcp/5517" help:"listen address in multiaddr format"`
}) {
	app := me.open(args.Config)
	app.RELAYMODE = []relay.RelayOpt{relay.OptHop}
	me.save(args.Config, app)
}

// Genconfig ...
func (me *T) Genconfig(args struct {
	Config string `default:".carb.json" help:"carb config file path"`
}) {
	app := new(app.T)
	app.SK = new(sk.T)
	app.SK.Gen()
	app.LISTEN = []string{"/ip4/0.0.0.0/tcp/5517", "/ip4/0.0.0.0/udp/5517/quic"}
	me.save(args.Config, app)
}

// Genkey ...
func (me *T) Genkey(args struct {
	Config string `default:".carb.json" help:"carb config file path"`
}) {
	app := me.open(args.Config)
	app.SK.Gen()
	me.save(args.Config, app)
}

// Node ...
func (me *T) Node(args struct {
	Config string `default:".carb.json" help:"carb config file path"`
}) {
	app := me.open(args.Config)
	app.Start()
	select {}
}

// Showid ...
func (me *T) Showid(args struct {
	Config string `default:".carb.json" help:"carb config file path"`
}) {
	app := me.open(args.Config)
	id, err := peer.IDFromPrivateKey(app.SK)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(id)
}

func (me *T) open(filename string) *app.T {
	app := new(app.T)
	conf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("[ReadConfig]:", err)
	}
	if err = json.Unmarshal(conf, app); err != nil {
		log.Fatalln("[JSON]:", err)
	}
	return app
}

func (me *T) save(filename string, app *app.T) {
	conf, err := json.MarshalIndent(app, "", "    ")
	if err != nil {
		log.Fatalln("[JSON]:", err)
	}
	if err = ioutil.WriteFile(filename, conf, 0600); err != nil {
		log.Fatalln("[WriteFile]:", err)
	}
}
