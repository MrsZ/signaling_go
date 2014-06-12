package signaling

import "net/http"
import "github.com/go-martini/martini"
import "github.com/martini-contrib/cors"


func App() *martini.ClassicMartini{
	m := martini.Classic()
	// Make a new Broker instance
	broker := NewBroker()

	m.Get("/", func() string {
		return "Sup"
	})
	// todo: ?broker injection
	m.Use(cors.Allow(&cors.Options{
		   AllowOrigins:     []string{"*"},
		   AllowMethods:     []string{"POST","OPTIONS"},
		   AllowHeaders:     []string{"Origin", "Content-type"},
	}))

	m.Post("/update/:room",
		func(resp http.ResponseWriter, req *http.Request, params martini.Params){
				UpdateHandler(resp, req, params, broker)
		})

	m.Options("/update/:room", OptionsHandler)

	m.Get("/stream/:room",
		func(resp http.ResponseWriter, req *http.Request, params martini.Params){
			ClientStream(resp, req, params, broker)
		})

	return m
}
