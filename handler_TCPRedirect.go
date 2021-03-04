package main

import (
	"io"
	"log"
	"net"
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
				go func() {
					reader := io.TeeReader(str, libConnDeadlineRefresher{conn, cfg.Idle})
					if n, err := io.Copy(conn, reader); err != nil {
						logger.Println("[Stream -> Conn]:", err, "{<n>}:", n)
					} else if cfg.PrintLog {
						logger.Println("[Stream -> Conn]:", n)
					}
				}()
				go func() {
					defer conn.Close()
					defer str.Close()
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
