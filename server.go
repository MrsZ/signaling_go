package main

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"github.com/martini-contrib/gorelic"
	"net/http"
	"signaling"
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
		gorelic.InitNewrelicAgent(settings.NewRelicKey, settings.NewRelicApp, true)
		martiniApp.Use(gorelic.Handler)
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
		GOMAXPROCS int
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
