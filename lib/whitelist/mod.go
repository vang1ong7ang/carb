package whitelist

import "github.com/libp2p/go-libp2p-core/peer"

// T ...
type T struct {
	data map[peer.ID]bool
}

// Add ...
func (me *T) Add(ids ...peer.ID) {
	me.data = make(map[peer.ID]bool)
	for _, v := range ids {
		me.data[v] = true
	}
}

// Query ...
func (me *T) Query(id peer.ID) bool {
	return me.data[id]
}
