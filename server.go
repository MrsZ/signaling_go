package main

import "net/http"
import "signaling"


func main() {
	martiniApp := signaling.App()
	http.ListenAndServe("0.0.0.0:8080", martiniApp)
}
