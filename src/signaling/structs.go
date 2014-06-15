package signaling


type Message struct {
	Id   string
	Data string
	Type string
	From string
	To string
	Room string
}


type RestrictedMsg struct {
	Type  string      `json:"type"`
	From  string `json:"from"`
	To string     `json:"to"`
}

const MaxMembers = 6
