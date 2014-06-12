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
	expect(t, len(roomWithGuy), 2)

}

