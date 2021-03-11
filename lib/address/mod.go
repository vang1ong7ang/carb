package address

import (
	"encoding/json"

	"github.com/multiformats/go-multiaddr"
)

// T ...
type T struct {
	multiaddr.Multiaddr
}

// UnmarshalJSON ...
func (me *T) UnmarshalJSON(p []byte) (err error) {
	var str string

	if err = json.Unmarshal(p, &str); err != nil {
		return
	}

	me.Multiaddr, err = multiaddr.NewMultiaddr(str)
	return
}
