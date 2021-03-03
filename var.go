package main

import (
	"os"

	"github.com/libp2p/go-libp2p-core/protocol"
)

var (
	cfg      config
	handlers = map[protocol.ID]interface{}{}
	clients  = map[protocol.ID]interface{}{}
	conffile = func() string {
		if file := os.ExpandEnv(`${CONFIG}`); len(file) > 0 {
			return file
		}
		return ".carb.json"
	}()
)
