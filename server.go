package main

import (
	"flag"
	"fmt"
	"github.com/martini-contrib/gorelic"
	"net/http"
	"os"
	"signaling"
)

func main() {
	addr := flag.String("bind", "0.0.0.0:8080", "Bind address ip:port")
	flag.Parse()
	martiniApp := signaling.App()

	newRelicKey := os.Getenv("NewRelicKey")
	if len(newRelicKey) > 0 {
		gorelic.InitNewrelicAgent(newRelicKey, "InstaMuteGo", true)
		martiniApp.Use(gorelic.Handler)
	}

	fmt.Printf("Start serving %s\n", *addr)
	http.ListenAndServe(*addr, martiniApp)
}
