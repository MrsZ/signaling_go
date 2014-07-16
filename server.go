package main

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"github.com/msoedov/signaling_go/signaling"
	"github.com/msoedov/signaling_go/newrelic"
	"net/http"
	"runtime"
)

// compile passing -ldflags "-X main.Build <build sha1>"
var Build string

func main() {
	filepath := flag.String("c", "default.ini", "Config file path")
	flag.Parse()
	martiniApp := signaling.App()
	settings := ReadConfig(*filepath).App
	if len(settings.NewRelicApp) > 0 {
		newrelic.InitNewrelicAgent(settings.NewRelicKey, settings.NewRelicApp, false, signaling.MembersBroker)
	}
	fmt.Printf("Using build: %s\n", Build)
	runtime.GOMAXPROCS(settings.GOMAXPROCS)
	fmt.Printf("Start serving %s\n", settings.Addr)
	http.ListenAndServe(settings.Addr, martiniApp)
}

type Settings struct {
	App struct {
		Addr        string
		NewRelicKey string
		NewRelicApp string
		GOMAXPROCS  int
	}
}

func ReadConfig(filename string) *Settings {
	var settings Settings
	err := gcfg.ReadFileInto(&settings, filename)
	if err != nil {
		panic(err)
	}
	return &settings
}
