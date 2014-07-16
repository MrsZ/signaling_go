package signaling_test

import (
	"bytes"
	"encoding/json"
	assert "github.com/msoedov/signaling_go/assert"
	. "github.com/msoedov/signaling_go/signaling"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MartiniResponseRecorder struct {
	// Deal with CloseNotify interface inside martini handlers
	// Pretty useful to send @CloseNotifyC <- true for close response
	// stream
	CloseNotifyC chan bool
	httptest.ResponseRecorder
}

func (self *MartiniResponseRecorder) CloseNotify() <-chan bool {
	return self.CloseNotifyC
}

func NewMartiniRecorder() *MartiniResponseRecorder {
	return &MartiniResponseRecorder{
		make(chan bool, 1),
		httptest.ResponseRecorder{
			HeaderMap: make(http.Header),
			Body:      new(bytes.Buffer),
			Code:      200,
		},
	}
}

// Parse buffer from format:
// event: @name\n
// data: @json nested structure\n
// \n
func ParsePayload(buf *bytes.Buffer) (name string, json_data map[string]string) {
	parts := strings.SplitN(buf.String(), "\n", -1)
	if len(parts) != 4 {
		panic("Malformed response")
	}
	name = strings.SplitN(parts[0], ":", -1)[1]
	name = strings.TrimSpace(name)
	var data string = strings.SplitN(parts[1], ":", 2)[1]
	data = strings.TrimSpace(data)
	json.Unmarshal([]byte(data), &json_data)
	return
}

func TestMainPage(t *testing.T) {
	response := httptest.NewRecorder()

	martiniApp := App()
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	assert.Equals(t, response.Code, http.StatusOK)
}

func TestEmptyUpdatePost(t *testing.T) {
	response := httptest.NewRecorder()

	martiniApp := App()
	request, err := http.NewRequest("POST", "/update/foo", strings.NewReader(""))
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	assert.Equals(t, response.Code, http.StatusBadRequest)
}

func TestMalformedUpdatePost(t *testing.T) {
	response := httptest.NewRecorder()

	martiniApp := App()
	request, err := http.NewRequest("POST", "/update/foo", strings.NewReader("foo-bar"))
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	assert.Equals(t, response.Code, http.StatusBadRequest)
}

func TestGoodUpdatePost(t *testing.T) {
	response := httptest.NewRecorder()

	martiniApp := App()
	body := string(`{"type":"invite",
					"from":"c2b0eb25",
					"to":"bf164412",
					"invite":"invite"}`)
	request, err := http.NewRequest("POST", "/update/foo", strings.NewReader(body))
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	assert.Equals(t, response.Code, http.StatusOK)
}

func TestUpdatePostTypeOnly(t *testing.T) {
	response := httptest.NewRecorder()

	martiniApp := App()
	body := []byte(`{"type":"invite",
					"invite":"invite"}`)
	request, err := http.NewRequest("POST", "/update/foo", strings.NewReader(string(body)))
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	assert.Equals(t, response.Code, http.StatusOK)
}

func TestUpdatePostNoType(t *testing.T) {
	response := httptest.NewRecorder()

	martiniApp := App()
	body := []byte(`{"from":"c2b0eb25",
					"to":"bf164412",
					"invite":"invite"}`)
	request, err := http.NewRequest("POST", "/update/foo", strings.NewReader(string(body)))
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	assert.Equals(t, response.Code, http.StatusBadRequest)
}

func TestUpdatePostNested(t *testing.T) {
	response := httptest.NewRecorder()

	martiniApp := App()
	body := []byte(`{"type":"offer",
					"from":"c2b0eb25",
					"to":"bf164412",
					"offer":{"type":"offer","sdp":"v=0\r\no=Mozilla-SIPUA-28.0 15773 0 IN IP4 0.0.0.0\r\ns"}}`)
	request, err := http.NewRequest("POST", "/update/foo", strings.NewReader(string(body)))
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	assert.Equals(t, response.Code, http.StatusOK)
}

func TestUpdatePostWrongPayloads(t *testing.T) {
	t.SkipNow()
	response := httptest.NewRecorder()

	martiniApp := App()
	body := []byte(`{"type":"answer",
					"from":"c2b0eb25",
					"to":"bf164412",
					"offer":{"type":"offer","sdp":"v=0\r\no=Mozilla-SIPUA\r\ns"}}`)
	request, err := http.NewRequest("POST", "/update/foo", strings.NewReader(string(body)))
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	assert.Equals(t, response.Code, http.StatusBadRequest)
}

func TestUpdateOptions(t *testing.T) {
	response := httptest.NewRecorder()

	martiniApp := App()
	request, err := http.NewRequest("OPTIONS", "/update/foo", strings.NewReader(""))
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	assert.Equals(t, response.Code, http.StatusOK)
}

func TestBrokerRoom(t *testing.T) {
	broker := NewBroker()
	room := broker.Room("foo")

	assert.Equals(t, len(room), 0)
	messageChan := make(chan *Message)
	room["SomeGuy"] = messageChan

	roomWithGuy := broker.Room("foo")
	assert.Equals(t, len(roomWithGuy), 1)

}

func TestBrokerRelease(t *testing.T) {
	broker := NewBroker()
	room := broker.Room("foo")

	assert.Equals(t, len(room), 0)
	messageChan := make(chan *Message)
	room["SomeGuy"] = messageChan

	broker.Release("foo", "SomeGuy")

	roomWithGuy := broker.Room("foo")
	assert.Equals(t, len(roomWithGuy), 0)

}

func TestStreamBasic(t *testing.T) {
	response := NewMartiniRecorder()

	martiniApp := App()
	//	GET and immediately close
	response.CloseNotifyC <- true
	request, err := http.NewRequest("GET", "/stream/foo", nil)
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	assert.Equals(t, response.Code, http.StatusOK)
	assert.Equals(t, response.HeaderMap["Access-Control-Allow-Methods"][0], "POST,OPTIONS")
	assert.Equals(t, response.HeaderMap["Content-Type"][0], "text/event-stream")
}

func TestStreamResponsePayload(t *testing.T) {
	response := NewMartiniRecorder()

	martiniApp := App()
	//	GET and immediately close
	response.CloseNotifyC <- true
	request, err := http.NewRequest("GET", "/stream/TestStreamResponsePayload", nil)
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	assert.Equals(t, response.Code, http.StatusOK)

	name, payloads := ParsePayload(response.Body)

	assert.Equals(t, name, "uid")

	assert.Equals(t, payloads["type"], "uid")
	assert.Equals(t, payloads["uid"], payloads["from"])
	//	present from and to
	_, ok := payloads["from"]
	assert.Assert(t, ok, "From field missed")
}

func TestFailuresStats(t *testing.T) {
	response := NewMartiniRecorder()

	martiniApp := App()

	_, failures := MembersBroker.GetStats()

	assert.Equals(t, failures, 0)
	request, err := http.NewRequest("POST", "/failure", nil)
	assert.Ok(t, err)

	martiniApp.ServeHTTP(response, request)

	assert.Equals(t, response.Code, http.StatusOK)

	_, failures = MembersBroker.GetStats()
	assert.Equals(t, failures, 1)
	//	ensure counter reset
	_, failures = MembersBroker.GetStats()
	assert.Equals(t, failures, 0)

}
