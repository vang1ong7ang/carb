package main

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

type libTargetSelector map[peer.ID]float64

func (me libTargetSelector) Select() peer.ID {
	for k := range me {
		return k
	}
	panic(me)
}
