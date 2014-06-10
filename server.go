package main

import (
	"github.com/go-martini/martini"
	"github.com/nu7hatch/gouuid"
)
import "log"
import "fmt"
import "encoding/json"
import "io"
import "bytes"
import "net/http"


type Message struct {
	Id   string
	Data string
	Type string
	From string
	To string
	Room string
}


func (self *Message) NewBuddy() *Message {
	self.Type = "newbuddy"
	data := map[string]string {	"uid": self.From, "from": self.From, "to": self.From, "type": "newbuddy"}
	self.Data = ToJsonString(&data)
	return self
}


func (self *Message) Uid() *Message {
	data := map[string]string {	"uid": self.From, "from": self.From, "to": self.From, "type": "uid"}
	self.Data = ToJsonString(&data)
	return self
}


func (self *Message) Dropped() *Message {
	self.Type = "dropped"
	data := map[string]string {	"from": self.From, "to": "", "type": "dropped"}
	self.Data = ToJsonString(&data)
	return self
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

func ToJsonString(info *map[string]string) string {
	var buf bytes.Buffer
	result, _ := json.Marshal(info)
	buf.Write(result)
	return buf.String()
}


func (self *RestrictedMsg) ReadJson(reader io.Reader) error {
	return ReadJson(reader, self)
}


type Broker struct {

	// room name -> client uid -> client chanel
	clients map[string]map[string]chan *Message

	// Channel into which messages are pushed to be broadcast out
	// to attahed clients.
	//
	messages chan *Message
}


func NewBroker() *Broker {
	b := &Broker{
		make(map[string]map[string]chan *Message),
		make(chan *Message),
	}
	return b
}


func pushMessage(msg *Message, broker *Broker){
	if msg.To != "" {
		//	Concrete destination
		room, ok := broker.clients[msg.Room]
		if !ok {
				log.Printf("No such room %s", msg.Room)
				return
		}
		client_channel, ok := room[msg.To]
		if !ok {
				log.Printf("No such patcipant %s in room %s", msg.To, msg.Room)
			return
		}
		client_channel <- msg
	} else {
		// Should be send for all in room
		for name, q := range broker.clients[msg.Room] {
			if msg.From != name {
				q <- msg
			}
		}
	}
	log.Printf("Broadcast message to %s clients", msg.Room)
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
	var roomName = params["room"]
	uid4, err := uuid.NewV4()
	if err != nil {
		http.Error(resp, "uid failed",
			http.StatusInternalServerError)
		return
	}
	var uid = uid4.String()
	room, ok := b.clients[roomName]
	if !ok {
			room = make(map[string] chan *Message)
			b.clients[roomName] = room
	}


	message := &Message{"", "", "uid", uid, "", roomName}

	var msg = message.Uid()
	fmt.Fprintf(resp, "event: %s\n", msg.Type)
	fmt.Fprintf(resp, "data: %s\n\n", msg.Data)
	f.Flush()


	pushMessage(message.NewBuddy(), b)
	room[uid] = messageChan
	// Remove this client from the map of attached clients
	// when `ClientStream` exits.
	defer func() {
		delete(room, uid)
		pushMessage(message.Dropped(), b)
		if len(room) == 0 {
			delete(b.clients, roomName)
			log.Println("Releasing room %s", roomName)
		}
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
	var roomName = params["room"]
	if err := ReadJson(buf, route); err != nil {
		http.Error(resp, "Bad Request", http.StatusBadRequest)
	} else {
		message := &Message{"", buf.String(), route.Type, route.From, route.To, roomName}
		pushMessage(message, b)
	}
	resp.WriteHeader(200)
}


func CorpMiddleware(resp http.ResponseWriter, req *http.Request){
	// todo: sophisticated middleware
	headers := resp.Header()
	headers.Set("Access-Control-Allow-Origin", "*")
}

func main() {
	m := martini.Classic()
	// Make a new Broker instance
	broker := NewBroker()
	m.Use(CorpMiddleware)

	m.Get("/", func() string {
		return "Sup"
	})
	m.Post("/update/:room", func(resp http.ResponseWriter, req *http.Request, params martini.Params){ UpdateHandler(resp, req, params, broker) })
	m.Options("/update/:room", func(resp http.ResponseWriter, req *http.Request, params martini.Params){ UpdateHandler(resp, req, params, broker) })
	m.Get("/stream/:room", func(resp http.ResponseWriter, req *http.Request, params martini.Params){ ClientStream(resp, req, params, broker) })

	http.ListenAndServe(":8080", m)
}
