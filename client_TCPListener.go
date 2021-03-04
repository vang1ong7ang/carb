package main

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

func init() {
	clients["TCPListener"] = func(
		cfg struct {
			Target        string
			ListenAddress string
			PrintLog      bool
		},
		logger *log.Logger,
	) func(hst host.Host) {
		return func(hst host.Host) {
			if listener, err := net.Listen("tcp", cfg.ListenAddress); err != nil {
				logger.Println("[Listen]:", err)
			} else if id, err := peer.IDB58Decode(cfg.Target); err != nil {
				logger.Println("[IDDecode]:", err)
			} else {
				defer listener.Close()
				for {
					if conn, err := listener.Accept(); err != nil {
						logger.Println("[Accept]:", err)
					} else if str, err := hst.NewStream(context.Background(), id, "TCPRedirect"); err != nil {
						logger.Println("[NewStream]:", err)
						conn.Close()
					} else {
						go func() {
							defer conn.Close()
							defer str.Close()
							if n, err := io.Copy(conn, str); err != nil {
								logger.Println("[Stream -> Conn]:", err, "{<n>}:", n)
							} else if cfg.PrintLog {
								logger.Println("[Stream -> Conn]:", n)
							}
						}()
						go func() {
							defer str.CloseWrite()
							if n, err := io.Copy(str, conn); err != nil {
								logger.Println("[Conn -> Stream]:", err, "{<n>}:", n)
							} else if cfg.PrintLog {
								logger.Println("[Conn -> Stream]:", n)
							}
						}()
					}
				}
			}
		}
	}
}
