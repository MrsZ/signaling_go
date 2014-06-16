package signaling


type Message struct {
	Id   string
	Data string
	Room string
	Meta
}


type Meta struct {
	Type  string  `json:"type"`
	From  string `json:"from"`
	To string     `json:"to"`
}

const MaxMembers = 6


func (self *Message) NewBuddy() *Message {
	self.Type = "newbuddy"
	data := map[string]string {"uid": self.From, "from": self.From, "to": "", "type": "newbuddy"}
	self.Data = ToJsonString(&data)
	return self
}


func (self *Message) Uid() *Message {
	data := map[string]string {"uid": self.From, "from": self.From, "to": self.From, "type": "uid"}
	self.Data = ToJsonString(&data)
	return self
}


func (self *Message) Dropped() *Message {
	self.Type = "dropped"
	data := map[string]string {"from": self.From, "to": "", "type": "dropped"}
	self.Data = ToJsonString(&data)
	return self
}


func (self *Message) Rejected() *Message {
	self.Type = "rejected"
	data := map[string]string {"from": self.From, "to": "", "type": "rejected",
		"message": "Room is full"}
	self.Data = ToJsonString(&data)
	return self
}

