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
func (me *T) Gen() error {
	sk, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)

	if err != nil {
		return err
	}

	me.PrivKey = sk
	return nil
}

// UnmarshalJSON ...
func (me *T) UnmarshalJSON(p []byte) error {
	str := new(string)
	err := json.Unmarshal(p, str)

	if err != nil {
		return err
	}

	dec, err := crypto.ConfigDecodeKey(*str)

	if err != nil {
		return err
	}

	sk, err := crypto.UnmarshalPrivateKey(dec)

	if err != nil {
		return err
	}

	me.PrivKey = sk
	return nil
}

// MarshalJSON ...
func (me *T) MarshalJSON() ([]byte, error) {
	enc, err := crypto.MarshalPrivateKey(me.PrivKey)

	if err != nil {
		return nil, err
	}

	return json.Marshal(crypto.ConfigEncodeKey(enc))
}
