package commander

import (
	"log"
	"os/exec"

	"github.com/libp2p/go-libp2p-core/host"
)

// T ...
type T struct {
	Name string
	Args []string
}

// Init ...
func (me *T) Init(hst host.Host, logger *log.Logger) {
	if err := exec.Command(me.Name, me.Args...).Run(); err != nil {
		logger.Println("[Error]:", err)
	}
}
