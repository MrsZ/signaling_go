package main

import "net/http"
import "github.com/go-martini/martini"
import "./signaling"


func main() {
	m := martini.Classic()
	// Make a new Broker instance
	broker := signaling.NewBroker()

	m.Get("/", func() string {
		return "Sup"
	})
	// todo: CORP middleware
	// todo: ?broker in context
	m.Post("/update/:room",
		func(resp http.ResponseWriter, req *http.Request, params martini.Params){
				signaling.UpdateHandler(resp, req, params, broker)
		})

	m.Options("/update/:room",
		func(resp http.ResponseWriter, req *http.Request, params martini.Params){
				signaling.UpdateHandler(resp, req, params, broker)
		})

	m.Get("/stream/:room",
		func(resp http.ResponseWriter, req *http.Request, params martini.Params){
			signaling.ClientStream(resp, req, params, broker)
		})

	//	todo: get addr from env
	http.ListenAndServe("0.0.0.0:8080", m)
}
