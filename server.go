package main

import (
	"github.com/go-martini/martini"
	"net/http"
)
import "log"
import "fmt"
import "encoding/json"
import "io"
import "bytes"


type Message struct {
	Id   string
	Data string
	Type string
	From string
	To string
}


type RestrictedMsg struct {
	Type  string      `json:"type"`
	From  string `json:"from"`
	To string     `json:"to"`
}

func ReadJson(from io.Reader, to interface{}) error {
	dec := json.NewDecoder(from)
	if err := dec.Decode(to); err != nil {
		log.Println(err)
		return err
	}
	return nil
}


func (self *RestrictedMsg) ReadJson(reader io.Reader) error {
	return ReadJson(reader, self)
}


type Broker struct {
	// Create a map of clients, the keys of the map are the channels
	// over which we can push messages to attached clients. (The values
	// are just booleans and are meaningless.)
	//
	clients map[chan *Message]bool

	// Channel into which new clients can be pushed
	//
	newClients chan chan *Message

	// Channel into which disconnected clients should be pushed
	//
	defunctClients chan chan *Message

	// Channel into which messages are pushed to be broadcast out
	// to attahed clients.
	//
	messages chan *Message
}


func NewBroker() *Broker {
	b := &Broker{
		make(map[chan *Message]bool),
		make(chan (chan *Message)),
		make(chan (chan *Message)),
		make(chan *Message),
	}
	return b
}

func (self *Broker) Start() {
	// Start a goroutine
	//
	go func() {
		// Loop endlessly
		//
		for {
			// Block until we receive from one of the
			// three following channels.
			select {
			case s := <-self.newClients:
				// There is a new client attached and we
				// want to start sending them messages.
				self.clients[s] = true
				log.Println("Added new client")
			case s := <-self.defunctClients:
				// A client has dettached and we want to
				// stop sending them messages.
				delete(self.clients, s)
				log.Println("Removed client")
			case msg := <-self.messages:
				// There is a new message to send. For each
				// attached client, push the new message
				// into the client's message channel.
				for s, _ := range self.clients {
					s <- msg
				}
				log.Printf("Broadcast message to %d clients", len(self.clients))
			}
		}
	}()
}


func ClientStream(resp http.ResponseWriter, req *http.Request, params martini.Params, b *Broker) {
	f, ok := resp.(http.Flusher)
	if !ok {
		http.Error(resp, "Streaming unsupported!",
			http.StatusInternalServerError)
		return
	}
	c, ok := resp.(http.CloseNotifier)
	if !ok {
		http.Error(resp, "close notification unsupported",
			http.StatusInternalServerError)
		return
	}
	// Create a new channel, over which the broker can
	// send this client messages.
	messageChan := make(chan *Message)
	// Add this client to the map of those that should
	// receive updates
	b.newClients <- messageChan
	// Remove this client from the map of attached clients
	// when `ClientStream` exits.
	defer func() {
		b.defunctClients <- messageChan
	}()
	headers := resp.Header()
	headers.Set("Content-Type", "text/event-stream")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Connection", "keep-alive")
	closer := c.CloseNotify()
	for {
		select {
		case msg := <-messageChan:
			if msg.Id != "" {
				fmt.Fprintf(resp, "id: %s\n", msg.Id)
			}
			fmt.Fprintf(resp, "event: %s\n", msg.Type)
			fmt.Fprintf(resp, "data: %s\n\n", msg.Data)
			f.Flush()
		case <-closer:
			log.Println("Closing connection")
			return
		}
	}
}

func UpdateHandler(resp http.ResponseWriter, req *http.Request, params martini.Params, b *Broker) {

	var route = new(RestrictedMsg)
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	if err := ReadJson(buf, route); err != nil {
		http.Error(resp, "Bad Request", http.StatusBadRequest)
	} else {
		message := &Message{"", buf.String(), route.Type, route.From, route.To}
		b.messages <- message
	}
	resp.WriteHeader(200)
}

//func CurryBroker(handler func(resp http.ResponseWriter, req *http.Request, b *Broker), broker *Broker) func(resp http.ResponseWriter, req *http.Request){
//	return func(resp http.ResponseWriter, req *http.Request){ handler(resp, req, broker) })
//}

func main() {
	m := martini.Classic()
	// Make a new Broker instance
	broker := NewBroker()

	m.Get("/", func() string {
		return "Sup"
	})
	m.Post("/update/:room", func(resp http.ResponseWriter, req *http.Request, params martini.Params){ UpdateHandler(resp, req, params, broker) })
	m.Get("/stream/:room", func(resp http.ResponseWriter, req *http.Request, params martini.Params){ ClientStream(resp, req, params, broker) })

	http.Handle("/sheet", m)
	// Start processing events
	broker.Start()
	http.ListenAndServe(":8080", m)
}
