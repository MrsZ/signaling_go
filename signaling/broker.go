package signaling

import (
	metrics "github.com/yvasiyarov/go-metrics"
	"log"
)

type Broker struct {

	// room name -> client uid -> client chanel
	clients   map[string]map[string]chan *Message
	connected metrics.Counter
	failures  metrics.Counter
}

func (self *Broker) Room(roomName string) map[string]chan *Message {
	room, ok := self.clients[roomName]
	if !ok {
		self.clients[roomName] = map[string]chan *Message{}
		return self.clients[roomName]
	} else {
		return room
	}
}

func (self *Broker) Release(roomName, member string) {
	room, ok := self.clients[roomName]
	if !ok {
		panic("Release non existed")
	}
	delete(room, member)
	if len(room) == 0 {
		delete(self.clients, roomName)
		log.Printf("Room %s was completly released", roomName)
	}
}

func NewBroker() *Broker {
	b := &Broker{
		make(map[string]map[string]chan *Message),
		metrics.NewCounter(),
		metrics.NewCounter(),
	}
	return b
}

func (broker *Broker) PushMessage(msg *Message) {
	if msg.To != "" {
		//Concrete message destination
		room, ok := broker.clients[msg.Room]
		if !ok {
			log.Printf("No such room %s", msg.Room)
			return
		}
		client_channel, ok := room[msg.To]
		if !ok {
			log.Printf("No such partcipant %s in room %s", msg.To, msg.Room)
			return
		}
		client_channel <- msg
	} else {
		//Global notification for all
		for name, q := range broker.clients[msg.Room] {
			if msg.From != name {
				q <- msg
			}
		}
	}
	log.Printf("Broadcast message %s to %s clients", msg.Type, msg.Room)
}

func (self *Broker) GetStats() (metrics.Counter, metrics.Counter) {
	return self.connected, self.failures
}
