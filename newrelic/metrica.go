package newrelic

type MembersMetrica struct {
	observable Observable
}

func (self *MembersMetrica) GetName() string {
	return "Active members"
}
func (self *MembersMetrica) GetUnits() string {
	return "Count"
}
func (self *MembersMetrica) GetValue() (float64, error) {
	memberCount, _ := self.observable.GetStats()
	return float64(memberCount.Count()), nil
}

type FailuresMetrica struct {
	observable Observable
}

func (self *FailuresMetrica) GetName() string {
	return "Failed to connect"
}
func (self *FailuresMetrica) GetUnits() string {
	return "Count"
}
func (self *FailuresMetrica) GetValue() (float64, error) {
	_, failures := self.observable.GetStats()
	defer func() { failures.Clear() }()
	return float64(failures.Count()), nil
}
