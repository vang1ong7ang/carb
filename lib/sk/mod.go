package sk

import (
	"crypto/rand"
	"encoding/json"

	"github.com/libp2p/go-libp2p-core/crypto"
)

// T ...
type T struct {
	crypto.PrivKey
}

// Gen ...
func (me *T) Gen() (err error) {
	me.PrivKey, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	return
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

// MarshalJSON ...
func (me T) MarshalJSON() (p []byte, err error) {
	if p, err = crypto.MarshalPrivateKey(me.PrivKey); err != nil {
		return
	}
	return json.Marshal(crypto.ConfigEncodeKey(p))
}
