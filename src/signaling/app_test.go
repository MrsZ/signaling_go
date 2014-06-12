package signaling

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"strings"
)


func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}


func TestMainPage(t *testing.T) {
	response := httptest.NewRecorder()

	martiniApp := App()
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	expect(t, response.Code, http.StatusOK)
}

func TestEmptyUpdatePost(t *testing.T) {
	response := httptest.NewRecorder()

	martiniApp := App()
	request, err := http.NewRequest("POST", "/update/foo", strings.NewReader(""))
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	expect(t, response.Code, http.StatusBadRequest)
}


func TestMalformedUpdatePost(t *testing.T) {
	response := httptest.NewRecorder()

	martiniApp := App()
	request, err := http.NewRequest("POST", "/update/foo", strings.NewReader("foo-bar"))
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	expect(t, response.Code, http.StatusBadRequest)
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

	expect(t, response.Code, http.StatusOK)
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

	expect(t, response.Code, http.StatusOK)
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

	expect(t, response.Code, http.StatusBadRequest)
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

	expect(t, response.Code, http.StatusOK)
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

	expect(t, response.Code, http.StatusBadRequest)
}


func TestUpdateOptions(t *testing.T) {
	response := httptest.NewRecorder()

	martiniApp := App()
	request, err := http.NewRequest("OPTIONS", "/update/foo", strings.NewReader(""))
	if err != nil {
		t.Fail()
	}
	martiniApp.ServeHTTP(response, request)

	expect(t, response.Code, http.StatusOK)
}


func TestBrokerRoom(t *testing.T){
	broker := NewBroker()
	room := broker.Room("foo")

	expect(t, len(room), 0)
	messageChan := make(chan *Message)
	room["SomeGuy"] = messageChan

	roomWithGuy := broker.Room("foo")
	expect(t, len(roomWithGuy), 1)

}


func TestBrokerRelease(t *testing.T){
	broker := NewBroker()
	room := broker.Room("foo")

	expect(t, len(room), 0)
	messageChan := make(chan *Message)
	room["SomeGuy"] = messageChan

	broker.Release("foo", "SomeGuy")

	roomWithGuy := broker.Room("foo")
	expect(t, len(roomWithGuy), 0)

}
