package summoners_test

import (
	assert "github.com/msoedov/signaling_go/assert"
	. "github.com/msoedov/signaling_go/summoners"
	"testing"
)

func TestNewNamesBasic(t *testing.T) {

	var name string = NewName(0)
	assert.Assert(t, name != "", "Expected non empty name")
}

func TestNewNamesCollision(t *testing.T) {

	var name_first string = NewName(0)
	var name_after string = NewName(1)
	assert.Assert(t, name_after != name_first, "Unaccepted name generation")
}

func TestNewNamesGeneration(t *testing.T) {
	for i := 0; i < 5; i++ {
		var name string = NewName(i)
		assert.Assert(t, name != "", "Expected non empty name")
	}
}
