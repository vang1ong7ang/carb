package main

import (
	"log"
	"os/exec"

	"github.com/libp2p/go-libp2p-core/host"
)

func init() {
	clients["KeyGenerator"] = func(
		cfg struct {
			Name string
			Args []string
		},
		logger *log.Logger,
	) func(hst host.Host) {
		return func(hst host.Host) {
			if err := exec.Command(cfg.Name, cfg.Args...).Run(); err != nil {
				logger.Println("[Exec]:", err)
			}
		}
	}
}
