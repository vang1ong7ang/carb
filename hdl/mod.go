package hdl

import (
	"carb/hdl/commander"
	"carb/hdl/forwardtcp"
	"carb/hdl/listentcp"
	"carb/hdl/peerlog"
	"carb/hdl/pinger"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path"
	"reflect"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
)

// T ...
type T struct {
	data map[protocol.ID]interface{}
}

// Init ...
func (me *T) Init() {
	me.data = make(map[protocol.ID]interface{})
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
		me.data[protocol.ID(path.Base(reflect.TypeOf(v).PkgPath()))] = v
	}
}

// GetConfig ...
func (me *T) GetConfig(id protocol.ID) (raw json.RawMessage, err error) {
	return json.MarshalIndent(me.data[id], "", "    ")
}

// Get ...
func (me *T) Get(id protocol.ID, cfg json.RawMessage, hst host.Host, wrt io.Writer) {
	var err error
	const flag = log.LstdFlags | log.Lmsgprefix
	val := reflect.ValueOf(me.data[id])
	ptr := reflect.New(val.Type())
	ptr.Elem().Set(val)
	if err = json.Unmarshal(cfg, ptr.Interface()); err != nil {
		log.Fatalln(err)
	}
	ptr.MethodByName("Init").Call([]reflect.Value{reflect.ValueOf(hst), reflect.ValueOf(log.New(wrt, fmt.Sprintf("{%s}: ", id), flag))})
}
