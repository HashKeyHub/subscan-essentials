package plugins

import (
	"flag"
	"fmt"
	"io/ioutil"
	"plugin"
	"reflect"
	"strings"
	
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/plugins/balance"
)

func init() {
	registerNative()
	registerStatic()
}

type PluginFactory Plugin

type PluginFactory2 interface {
	PluginFactory
	InitHttp2(e *bm.Engine)
}

var RegisteredPlugins = make(map[string]PluginFactory2)

func Register(name string, f interface{}) {
	if f == nil {
		return
	}

	if _, ok := RegisteredPlugins[name]; ok {
		return
	}

	RegisteredPlugins[name] = f.(PluginFactory2)

}

func List() []string {
	plugins := make([]string, 0, len(RegisteredPlugins))
	for name := range RegisteredPlugins {
		plugins = append(plugins, name)
	}
	return plugins
}

func registerNative() {
	Register("account", balance.New())
}

func registerStatic() {
	flag.Parse()
	pluginsDir := fmt.Sprintf("%s/plugins", flag.Lookup("conf").Value)
	files, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		p, err := plugin.Open(fmt.Sprintf("%s/%s", pluginsDir, file.Name()))
		if err != nil {
			panic(err)
		}
		if file.IsDir() {
			return
		}
		pluginName := strings.Split(file.Name(), ".")[0]
		f, err := p.Lookup("New")
		if err != nil {
			panic(err)
		}
		Register(pluginName, reflect.ValueOf(f).Call(nil)[0].Interface())
	}
}
