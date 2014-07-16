package newrelic_test

import (
	assert "github.com/msoedov/signaling_go/assert"
	. "github.com/msoedov/signaling_go/newrelic"
	"testing"
)

type TestObservable struct {
	normal, failures int
}

func (self *TestObservable) GetStats() (int, int) {
	return self.normal, self.failures
}

func TestMembersMetricaBasic(t *testing.T) {
	observable := new(TestObservable)
	m := &MembersMetrica{observable}
	val, err := m.GetValue()
	assert.Equals(t, val, 0.0)
	assert.Ok(t, err)
}

func TestFailuresMetricaBasic(t *testing.T) {
	observable := new(TestObservable)
	m := &FailuresMetrica{observable}
	val, err := m.GetValue()
	assert.Equals(t, val, 0.0)
	assert.Ok(t, err)
}
