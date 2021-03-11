package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	relay "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
)

type config struct {
	SK        string
	LISTEN    []string
	RELAYMODE []relay.RelayOpt
	PEERS     []struct {
		ID      peer.ID
		Address string
	}
	PROTOCOLS []struct {
		ID     protocol.ID
		CONFIG json.RawMessage
	}
	CLIENTS []struct {
		ID     protocol.ID
		CONFIG json.RawMessage
	}
}

func (me *config) parse() {
	if content, err := ioutil.ReadFile(conffile); err != nil {
		log.Println("[ReadConfigFile]:", err)
	} else if err := json.Unmarshal(content, me); err != nil {
		log.Fatalln("[JSON]:", err)
	}
}

func (me *config) sk() crypto.PrivKey {
	if bsk, err := crypto.ConfigDecodeKey(me.SK); err != nil {
		log.Fatalln("[ConfigDecode]:", err)
	} else if sk, err := crypto.UnmarshalPrivateKey(bsk); err != nil {
		log.Fatalln("[UnmarshalPrivateKey]:", err)
	} else {
		return sk
	}
	return nil
}
