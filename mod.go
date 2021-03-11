package main

import (
	"carb/app"
	"carb/hdl"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/libp2p/go-libp2p-core/protocol"
)

// T ...
type T struct{}

// Node ...
func (me *T) Node(args struct {
	Config string `default:".carb.json" help:"carb config file path"`
}) {
	var app app.T
	var conf []byte
	var err error
	if conf, err = ioutil.ReadFile(args.Config); err != nil {
		log.Fatalln("[ReadConfig]:", err)
	}
	if err = json.Unmarshal(conf, &app); err != nil {
		log.Fatalln("[JSON]:", err)
	}
	app.Start()
	select {}
}

// Genconfig ...
func (me *T) Genconfig(args struct {
	Config string `default:".carb.json" help:"carb config file path"`
}) {
	var app app.T
	var conf []byte
	var err error
	app.SK.Gen()
	app.LISTEN = []string{"/ip4/0.0.0.0/tcp/5517", "/ip4/0.0.0.0/udp/5517/quic"}
	if conf, err = json.MarshalIndent(app, "", "    "); err != nil {
		log.Fatalln("[JSON]:", err)
	}
	if err = ioutil.WriteFile(args.Config, conf, 0600); err != nil {
		log.Fatalln("[WriteFile]:", err)
	}
}

// Addhandler ...
func (me *T) Addhandler(args struct {
	Config string `default:".carb.json" help:"carb config file path"`
	ID     string `default:"pinger" help:"handler protocol id"`
}) {
	var app app.T
	var hdlmgt hdl.T
	var conf []byte

	var err error
	if conf, err = ioutil.ReadFile(args.Config); err != nil {
		log.Fatalln("[ReadConfig]:", err)
	}
	if err = json.Unmarshal(conf, &app); err != nil {
		log.Fatalln("[JSON]:", err)
	}
	hdlmgt.Init()

	var handler struct {
		ENABLE bool
		LOG    bool
		ID     protocol.ID
		CONFIG json.RawMessage
	}

	if handler.CONFIG, err = hdlmgt.GetConfig(protocol.ID(args.ID)); err != nil {
		log.Fatalln(err)
	}

	handler.ENABLE = true
	handler.ID = protocol.ID(args.ID)
	app.Handlers = append(app.Handlers, handler)

	if conf, err = json.MarshalIndent(app, "", "    "); err != nil {
		log.Fatalln("[JSON]:", err)
	}
	if err = ioutil.WriteFile(args.Config, conf, 0600); err != nil {
		log.Fatalln("[WriteFile]:", err)
	}
}
