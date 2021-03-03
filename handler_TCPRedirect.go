package main

import (
	"io"
	"log"
	"net"

	"github.com/libp2p/go-libp2p-core/network"
)

func init() {
	handlers["TCPRedirect"] = func(
		cfg struct {
			Target   string
			PrintLog bool
		},
		logger *log.Logger,
	) func(str network.Stream) {
		return func(str network.Stream) {
			// TODO: ADD ID AUTH HERE
			if conn, err := net.Dial("tcp", cfg.Target); err != nil {
				logger.Println("[DIAL]:", err)
				str.Reset()
			} else {
				go func() {
					if n, err := io.Copy(conn, str); err != nil {
						logger.Println("[OUT]:", err)
					} else if cfg.PrintLog {
						logger.Println("[WriteOUT]:", n)
					}
				}()
				go func() {
					defer conn.Close()
					defer str.Close()
					if n, err := io.Copy(str, conn); err != nil {
						logger.Println("[IN]:", err)
					} else if cfg.PrintLog {
						logger.Println("[WriteIN]:", n)
					}
				}()
			}
		}
	}
}
