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

type Runnable interface {
	RunString(entry string) api.EvalResult
}

type RemoteRunner struct {
	client *http.Client
}

func NewRemoteRunner() *RemoteRunner {
	return &RemoteRunner{
		client: new(http.Client),
	}
}

func (r *RemoteRunner) RunString(entry string) api.EvalResult {
	body := bytes.NewBufferString(entry)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:9119/v1", body)
	res := api.EvalResult{}
	if err != nil {
		res.Error = err.Error()
		return res
	}
	ap := ActionParams{
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

type LocalRunner struct {
	vm     *goja.Runtime
	config api.Config
}

func NewLocalRunner(cfg api.Config) *LocalRunner {
	vm := goja.New()
	local := &LocalRunner{vm: vm, config: cfg}
	InitBuiltins(vm)
	// init all plugins
	for _, each := range cfg.Plugins {
		ns := each.Namespace()
		log.Println("init plugin", ns)
		vm.Set(ns, each)
		if err := each.Init(vm); err != nil {
			log.Fatal(err)
		}
	}
	local.initInternals()
	if cfg.Setup != nil {
		log.Println("custom setting up Javascript virtual machine")
		if err := cfg.Setup(vm); err != nil {
			log.Fatal(err)
		}
	}
	rand.Seed(time.Now().UnixNano())
	return local
}

func (r *LocalRunner) RunString(entry string) api.EvalResult {
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
			res.Data = Print(val)
			res.Datatype = fmt.Sprintf("%T", val)
		}
	}
	return res
}

func (r *LocalRunner) initInternals() {
	r.vm.Set("_plugins", r.pluginInfo)
	r.vm.Set("_variables", r.globalVariables)
	r.vm.Set("_toggledebug", r.toggleDebug)
	r.vm.Set("_showhelp", r.showHelp)
	r.vm.Set("_methods", r.showMethods)
	r.vm.Set("_browse", r.browseObject)
	r.vm.Set("_markdowninspect", r.markdownInspect)
}

func (r *LocalRunner) browseObject(v any) any {
	if v == nil {
		return ""
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

func (r *LocalRunner) showMethods(v any) PlainText {
	if v == nil {
		return ""
	}
	rt := reflect.TypeOf(v)
	b := new(strings.Builder)
	printMethods(b, rt)
	return PlainText(b.String())
}

func (r *LocalRunner) showHelp() string {
	return "no help yet"
}

func (r *LocalRunner) toggleDebug() {
	if api.Debug {
		api.Debug = false
		log.Println("verbose log disabled")
	} else {
		api.Debug = true
		log.Println("verbose log enabled")
	}
}

func (r *LocalRunner) pluginInfo() (list []string) {
	for _, each := range r.config.Plugins {
		list = append(list, each.Namespace())
	}
	return
}

func (r *LocalRunner) globalVariables() (filtered []string) {
	for _, each := range r.vm.GlobalObject().Keys() {
		// skip internal var and funcs
		if strings.HasPrefix(each, "_") {
			continue
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

func (r *LocalRunner) Include(path string) api.EvalResult {
	data, err := os.ReadFile(path)
	if err != nil {
		return api.EvalResult{
			Error: err.Error(),
		}
	}
	return r.RunString(string(data))
}

func (r *LocalRunner) markdownInspect(v any) any {
	return PlainText(PrintMarkdown(v))
}

// Open calls the OS default program for uri
func open(uri string) error {
	var run string
	switch {
	case "windows" == runtime.GOOS:
		run = "start"
	case "darwin" == runtime.GOOS:
		run = "open"
	case "linux" == runtime.GOOS:
		run = "xdg-open"
	default:
		return fmt.Errorf("Unable to open uri:%v on:%v", uri, runtime.GOOS)
	}
	return exec.Command(run, uri).Start()
}
