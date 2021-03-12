package commander

import (
	"carb/lib/logger"
	"os/exec"

	"github.com/libp2p/go-libp2p-core/host"
)

// T ...
type T struct {
	Name string
	Args []string
}

// Init ...
func (me *T) Init(hst host.Host, log *logger.T) {
	err := exec.Command(me.Name, me.Args...).Run()
	if err != nil {
		log.Log(1, "exec", map[string]interface{}{"err": err})
	}
}
