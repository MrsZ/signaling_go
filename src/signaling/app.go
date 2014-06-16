package signaling

import "github.com/go-martini/martini"
import "github.com/martini-contrib/cors"

func App() *martini.ClassicMartini {
	m := martini.Classic()
	// Make a new Broker instance
	broker := NewBroker()
	m.Map(broker)

	m.Get("/", func() string {
		return "Sup"
	})
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-type"},
	}))

	m.Post("/update/:room", UpdateHandler)

	m.Options("/update/:room", OptionsHandler)

	m.Get("/stream/:room", ClientStream)

	return m
}
