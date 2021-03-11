package sk

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p-core/crypto"
)

// T ...
type T struct {
	crypto.PrivKey
}

// UnmarshalJSON ...
func (me *T) UnmarshalJSON(p []byte) (err error) {
	var str string

	if err = json.Unmarshal(p, &str); err != nil {
		return
	}

	if p, err = crypto.ConfigDecodeKey(str); err != nil {
		return
	}

	me.PrivKey, err = crypto.UnmarshalPrivateKey(p)
	return
}
