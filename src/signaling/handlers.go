package signaling

import (
	"github.com/go-martini/martini"
	"github.com/nu7hatch/gouuid"
)

import "log"
import "fmt"
import "encoding/json"
import "bytes"
import "net/http"


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
	// todo: some human friendly names
	uid4, err := uuid.NewV4()
	if err != nil {
		http.Error(resp, "uid failed",
			http.StatusInternalServerError)
		return
	}
	var uid = uid4.String()
	// todo: add max members checking
	room := b.Room(roomName)

	headers := resp.Header()
	headers.Set("Content-Type", "text/event-stream")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Connection", "keep-alive")
	f.Flush()
	closer := c.CloseNotify()

	message := &Message{"", "", "uid", uid, "", roomName}

	var msg = message.Uid()
	// todo: add local closure for send
	fmt.Fprintf(resp, "event: %s\n", msg.Type)
	fmt.Fprintf(resp, "data: %s\n", msg.Data)
	f.Flush()


	b.PushMessage(message.NewBuddy())
	room[uid] = messageChan
	// Remove this client from the map of attached clients
	// when `ClientStream` exits.
	defer func() {
		b.Release(roomName, uid)
		b.PushMessage(message.Dropped())
	}()

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
	buf := new(bytes.Buffer)
	bytes_read, _ := buf.ReadFrom(req.Body)

	var roomName = params["room"]
	log.Printf("Readed %d bytes from response", bytes_read)

	var data map[string]string
	json.Unmarshal(buf.Bytes(), &data)

	if len(data) == 0 || data["type"] == ""{
		http.Error(resp, "Bad Request", http.StatusBadRequest)
		return
	}
	log.Println("Ok", data)

	message := &Message{"", buf.String(), data["type"], data["from"], data["to"], roomName}
	b.PushMessage(message)

	resp.WriteHeader(200)
}


func OptionsHandler(resp http.ResponseWriter, req *http.Request){
	resp.WriteHeader(200)
}
