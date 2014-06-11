package signaling

import "net/http"
import "github.com/go-martini/martini"


func App() *martini.ClassicMartini{
	m := martini.Classic()
	// Make a new Broker instance
	broker := NewBroker()

	m.Get("/", func() string {
		return "Sup"
	})
	// todo: CORP middleware
	// todo: ?broker in context
	m.Post("/update/:room",
		func(resp http.ResponseWriter, req *http.Request, params martini.Params){
				UpdateHandler(resp, req, params, broker)
		})

	m.Options("/update/:room",
		func(resp http.ResponseWriter, req *http.Request, params martini.Params){
				UpdateHandler(resp, req, params, broker)
		})

	m.Get("/stream/:room",
		func(resp http.ResponseWriter, req *http.Request, params martini.Params){
			ClientStream(resp, req, params, broker)
		})

	return m
}
