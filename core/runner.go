package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/emicklei/simone/api"
)

type runnable interface {
	RunString(entry string) api.EvalResult
}

type remoteRunner struct {
	client *http.Client
}

func newRemoteRunner() *remoteRunner {
	return &remoteRunner{
		client: new(http.Client),
	}
}

func (r *remoteRunner) RunString(entry string) api.EvalResult {
	body := bytes.NewBufferString(entry)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:9119/v1", body)
	res := api.EvalResult{}
	if err != nil {
		res.Error = err.Error()
		return res
	}
	ap := actionParams{
		Debug:  true,
		Line:   "1",
		File:   "cli",
		Action: "eval",
	}
	ap.Inject(req.URL)
	resp, err := r.client.Do(req)
	if err != nil {
		res.Error = err.Error()
		return res
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res.Error = err.Error()
		return res
	}
	if err := json.Unmarshal(result, &res); err != nil {
		res.Error = err.Error()
	}
	return res
}

type localRunner struct {
	vm     *goja.Runtime
	config api.Config
}

func newLocalRunner(cfg api.Config) *localRunner {
	vm := goja.New()
	local := &localRunner{vm: vm, config: cfg}
	initBuiltins(vm)
	var ctx api.PluginContext = local
	// init all plugins
	for _, each := range cfg.Plugins {
		ns := each.Namespace()
		log.Println("init plugin", ns)
		vm.Set(ns, each)
		if err := each.Init(ctx); err != nil {
			log.Fatal(err)
		}
	}
	local.initInternals()
	if cfg.Setup != nil {
		log.Println("custom setting up Javascript virtual machine")
		if err := cfg.Setup(ctx); err != nil {
			log.Fatal(err)
		}
	}
	rand.Seed(time.Now().UnixNano())
	return local
}

// Set implements api.PluginContext
func (r *localRunner) Set(name string, value any) error {
	return r.vm.Set(name, value)
}

func (r *localRunner) RunString(entry string) api.EvalResult {
	res := api.EvalResult{}
	result, err := r.vm.RunString(entry)
	if err != nil {
		// syntax error
		res.Error = err.Error()
	} else {
		val := result.Export()
		if err, ok := val.(error); ok {
			// evaluation error
			res.Error = err.Error()
		} else {
			// no error
			res.RawData = val
			res.Datatype = fmt.Sprintf("%T", val)
		}
	}
	return res
}

func (r *localRunner) initInternals() {
	r.vm.Set("_plugins", r.pluginInfo)
	r.vm.Set("_variables", r.globalVariables)
	r.vm.Set("_toggledebug", r.toggleDebug)
	r.vm.Set("_showhelp", r.showHelp)
	r.vm.Set("_methods", r.showMethods)
	r.vm.Set("_browse", r.browseObject)
	r.vm.Set("_markdowninspect", r.markdownInspect)
	r.vm.Set("_login", r.handleLogin)
}

func (r *localRunner) browseObject(v any) any {
	if v == nil {
		return "null"
	}
	// store value in temporary variable TODO cleanup?
	key := "_" + randSeq(10) // make it internal such that :v will not show it
	r.vm.Set(key, v)
	return open(fmt.Sprintf("http://%s/v1?action=browse&source=%s", r.config.HostPort(), key))
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (r *localRunner) handleLogin(plugin any, username, password string) error {
	if handler, ok := plugin.(api.LoginHandler); ok {
		return handler.Login(username, password)
	}
	return fmt.Errorf("%v cannot handle login", plugin)
}

func (r *localRunner) showMethods(v any) PlainText {
	if v == nil {
		return ""
	}
	rt := reflect.TypeOf(v)
	b := new(strings.Builder)
	printMethods(b, rt)
	return PlainText(b.String())
}

func (r *localRunner) showHelp() string {
	return "no help yet"
}

func (r *localRunner) toggleDebug() {
	if api.Debug {
		api.Debug = false
		log.Println("verbose log disabled")
	} else {
		api.Debug = true
		log.Println("verbose log enabled")
	}
}

func (r *localRunner) pluginInfo() (list []string) {
	for _, each := range r.config.Plugins {
		list = append(list, each.Namespace())
	}
	return
}

func (r *localRunner) globalVariables() (filtered []string) {
	for _, each := range r.vm.GlobalObject().Keys() {
		// skip internal var and funcs (unless debugging)
		if !api.Debug {
			if strings.HasPrefix(each, "_") {
				continue
			}
		}
		v := r.vm.Get(each)
		vt := v.ExportType()
		if vt.Kind() == reflect.Func {
			continue
		}
		// https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
		var plugin api.Plugin
		pluginType := reflect.TypeOf(&plugin).Elem()
		if vt.Implements(pluginType) {
			continue
		}
		filtered = append(filtered, each)
	}
	return
}

func (r *localRunner) Include(path string) api.EvalResult {
	data, err := os.ReadFile(path)
	if err != nil {
		return api.EvalResult{
			Error: err.Error(),
		}
	}
	log.Printf("\033[1;34mrunning %s\033[0m\n", path)
	return r.RunString(string(data))
}

func (r *localRunner) markdownInspect(v any) string {
	return PrintMarkdown(v)
}

// Open calls the OS default program for uri
func open(uri string) error {
	switch {
	case "windows" == runtime.GOOS:
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", uri).Start()
	case "darwin" == runtime.GOOS:
		return exec.Command("open", uri).Start()
	case "linux" == runtime.GOOS:
		return exec.Command("xdg-open", uri).Start()
	default:
		return fmt.Errorf("Unable to open uri:%v on:%v", uri, runtime.GOOS)
	}
}
