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
func (me *T) UnmarshalJSON(p []byte) error {
	str := new(string)
	err := json.Unmarshal(p, &str)

	if err != nil {
		return err
	}

	addr, err := multiaddr.NewMultiaddr(*str)

	if err != nil {
		return err
	}

	me.Multiaddr = addr
	return nil
}
