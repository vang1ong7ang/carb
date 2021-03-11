package main

import (
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

func init() {
	handlers["TCPRedirect"] = func(
		cfg struct {
			Permissions map[peer.ID]bool
			Target      string
			Idle        time.Duration
			PrintLog    bool
		},
		logger *log.Logger,
	) func(str network.Stream) {
		return func(str network.Stream) {
			if cfg.Permissions[str.Conn().RemotePeer()] == false {
				str.Close()
				logger.Println("[PermissionDenied]:", str.Conn().RemotePeer())
				return
			}
			if conn, err := net.DialTimeout("tcp", cfg.Target, time.Second); err != nil {
				logger.Println("[DIAL]:", err)
				str.Reset()
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
