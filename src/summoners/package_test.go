package summoners

import (
	"testing"
)

func TestNewNamesBasic(t *testing.T) {

	var name string = NewName(0)
	if name == "" {
		t.Errorf("Expected non empty name")
	}
}

func TestNewNamesCollision(t *testing.T) {

	var name_first string = NewName(0)
	var name_after string = NewName(1)
	if name_after == name_first {
		t.Errorf("Unaccepted name generation")
	}
}

func TestNewNamesGeneration(t *testing.T) {
	for i := 0; i < 5; i++ {
		var name string = NewName(i)
		if name == "" {
			t.Errorf("Expected non empty name")
		}
	}
}
