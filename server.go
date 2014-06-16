package main

import (
	"flag"
	"fmt"
	"net/http"
	"signaling"
)

func main() {
	addr := flag.String("bind", "0.0.0.0:8080", "Bind address ip:port")
	flag.Parse()
	martiniApp := signaling.App()
	fmt.Printf("Start serving %s\n", *addr)
	http.ListenAndServe(*addr, martiniApp)
}
