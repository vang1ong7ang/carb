package ext

import (
	"carb/ext/commander"
	"carb/ext/forwardtcp"
	"carb/ext/listentcp"
	"carb/ext/peerlog"
	"carb/ext/pinger"
	"carb/lib/logger"
	"encoding/json"
	"log"
	"path"
	"reflect"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
)

// T ...
type T struct {
	data map[string]interface{}
}

// Init ...
func (me *T) Init() {
	me.data = make(map[string]interface{})
	for _, v := range []interface{}{
		forwardtcp.T{
			To:          "127.0.0.1:1080",
			IdleTimeout: 10 * time.Second,
		},
		pinger.T{
			TimeInterval: time.Second,
		},
		commander.T{},
		peerlog.T{
			TimeInterval: time.Minute,
		},
		listentcp.T{
			Listen:      "127.0.0.1:1080",
			IdleTimeout: 10 * time.Second,
		},
	} {
		me.data[path.Base(reflect.TypeOf(v).PkgPath())] = v
	}
}

// GetConfig ...
func (me *T) GetConfig(id string) (raw json.RawMessage, err error) {
	return json.MarshalIndent(me.data[id], "", "    ")
}

// Get ...
func (me *T) Get(id string, cfg json.RawMessage, hst host.Host, lvl int) {
	val := reflect.ValueOf(me.data[id])
	ptr := reflect.New(val.Type())
	ptr.Elem().Set(val)
	err := json.Unmarshal(cfg, ptr.Interface())

	if err != nil {
		log.Fatalln(err)
	}

	logger := new(logger.T)
	logger.Init(id)
	logger.SetLevel(lvl)

	ptr.MethodByName("Init").Call([]reflect.Value{reflect.ValueOf(hst), reflect.ValueOf(logger)})
}
