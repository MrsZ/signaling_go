package newrelic_test

import (
	assert "github.com/msoedov/signaling_go/assert"
	. "github.com/msoedov/signaling_go/newrelic"
	metrics "github.com/yvasiyarov/go-metrics"
	"testing"
)

type TestObservable struct {
	normal   metrics.Counter
	failures metrics.Counter
}

func (self *TestObservable) GetStats() (metrics.Counter, metrics.Counter) {
	return self.normal, self.failures
}

func TestMembersMetricaBasic(t *testing.T) {
	observable := TestObservable{metrics.NewCounter(), metrics.NewCounter()}
	m := &MembersMetrica{&observable}
	val, err := m.GetValue()
	assert.Equals(t, val, 0.0)
	assert.Ok(t, err)
}

func TestFailuresMetricaBasic(t *testing.T) {
	observable := TestObservable{metrics.NewCounter(), metrics.NewCounter()}
	m := &FailuresMetrica{&observable}
	val, err := m.GetValue()
	assert.Equals(t, val, 0.0)
	assert.Ok(t, err)
}

func TestFailuresMetricaStats(t *testing.T) {
	failures := metrics.NewCounter()
	observable := TestObservable{metrics.NewCounter(), failures}
	failures.Inc(13)
	m := &FailuresMetrica{&observable}
	val, err := m.GetValue()
	assert.Equals(t, val, 13.0)
	assert.Ok(t, err)
	assert.Equals(t, failures.Count(), int64(0))
	// after all shoul be 0
	val, err = m.GetValue()
	assert.Equals(t, val, 0.0)
	assert.Ok(t, err)
}
