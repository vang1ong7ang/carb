package main

import (
	"context"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
)

func init() {
	clients["TCPListener"] = func(
		cfg struct {
			Target        peer.ID
			ListenAddress string
			Idle          time.Duration
			PrintLog      bool
		},
		logger *log.Logger,
	) func(hst host.Host) {
		return func(hst host.Host) {
			if listener, err := net.Listen("tcp", cfg.ListenAddress); err != nil {
				logger.Println("[Listen]:", err)
			} else {
				defer listener.Close()
				for {
					if conn, err := listener.Accept(); err != nil {
						logger.Println("[Accept]:", err)
					} else if str, err := hst.NewStream(context.Background(), cfg.Target, "TCPRedirect"); err != nil {
						logger.Println("[NewStream]:", err)
						conn.Close()
					} else {
						con := &libConnDeadlineRefresher{conn, cfg.Idle}
						con.Refresh()
						var wg sync.WaitGroup
						wg.Add(2)
						defer con.Close()
						defer str.Close()
						defer wg.Wait()

						go func() {
							defer wg.Done()
							defer str.CloseRead()
							if n, err := io.Copy(con, str); err != nil {
								logger.Println("[Str -> Con]:", err, "{<n>}:", n)
							} else if cfg.PrintLog {
								logger.Println("[Str -> Con]:", n)
							}
						}()

						go func() {
							defer wg.Done()
							defer str.CloseWrite()
							if n, err := io.Copy(str, con); err != nil {
								logger.Println("[Con -> Str]:", err, "{<n>}:", n)
							} else if cfg.PrintLog {
								logger.Println("[Con -> Str]:", n)
							}
						}()
					}
				}
			}
		}
	}
}
