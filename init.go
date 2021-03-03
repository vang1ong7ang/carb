package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/libp2p/go-libp2p-core/crypto"
)

func init() {
	cfg.parse()

	cmd := flag.String("Command", "", "command")
	gk := flag.Bool("GenerateKey", false, "...")
	ok := flag.Bool("OverwriteKey", false, "...")

	flag.Parse()

	switch *cmd {
	case "config":
		if *gk && (len(cfg.SK) == 0 || *ok) {
			if sk, pk, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader); err != nil {
				log.Fatalln(err)
			} else if bsk, err := crypto.MarshalPrivateKey(sk); err != nil {
				log.Fatalln(err)
			} else if bpk, err := crypto.MarshalPublicKey(pk); err != nil {
				log.Fatalln(err)
			} else {
				log.Println("[GenerateKey]:", "{PublicKey}:", crypto.ConfigEncodeKey(bpk))
				cfg.SK = crypto.ConfigEncodeKey(bsk)
			}
		}

		if content, err := json.Marshal(cfg); err != nil {
			log.Println("[JSONEncode]:", err)
		} else if err := ioutil.WriteFile(conffile, content, 0600); err != nil {
			log.Println("[WriteFile]:", err)
		}
		os.Exit(0)
	default:
	}
}
