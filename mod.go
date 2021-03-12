package main

import (
	"carb/app"
	"carb/ext"
	"carb/lib/sk"
	"encoding/json"
	"io/ioutil"
	"log"
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
	app.SK = new(sk.T)
	app.SK.Gen()
	app.LISTEN = []string{"/ip4/0.0.0.0/tcp/5517", "/ip4/0.0.0.0/udp/5517/quic"}
	if conf, err = json.MarshalIndent(app, "", "    "); err != nil {
		log.Fatalln("[JSON]:", err)
	}
	if err = ioutil.WriteFile(args.Config, conf, 0600); err != nil {
		log.Fatalln("[WriteFile]:", err)
	}
}

// Addext ...
func (me *T) Addext(args struct {
	Config string `default:".carb.json" help:"carb config file path"`
	ID     string `default:"pinger" help:"handler protocol id"`
}) {
	var app app.T
	var hdlmgt ext.T
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
		ENABLE   bool
		LOGLEVEL int
		ID       string
		CONFIG   json.RawMessage
	}

	if handler.CONFIG, err = hdlmgt.GetConfig(args.ID); err != nil {
		log.Fatalln(err)
	}

	handler.ENABLE = true
	handler.ID = args.ID
	app.Handlers = append(app.Handlers, handler)

	if conf, err = json.MarshalIndent(app, "", "    "); err != nil {
		log.Fatalln("[JSON]:", err)
	}
	if err = ioutil.WriteFile(args.Config, conf, 0600); err != nil {
		log.Fatalln("[WriteFile]:", err)
	}
}
