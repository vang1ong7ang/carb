package logger

import (
	"fmt"
	"log"
	"os"
)

// T ...
type T struct {
	val *log.Logger
	lvl int
}

// Init ...
func (me *T) Init(prefix string) {
	me.val = log.New(os.Stderr, fmt.Sprintf("{{ %s }}: ", prefix), log.LstdFlags|log.Lmsgprefix)
}

// SetLevel ...
func (me *T) SetLevel(lvl int) {
	me.lvl = lvl
}

// Log ...
func (me *T) Log(lvl int, k string, v map[string]interface{}) {
	if lvl > me.lvl {
		return
	}

	me.val.Printf("[[ %s ]]: %v", k, v)
}
