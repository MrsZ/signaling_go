package newrelic

type MembersMetrica struct {
	Observable Observable
}

func (self *MembersMetrica) GetName() string {
	return "Active members"
}
func (self *MembersMetrica) GetUnits() string {
	return "Count"
}
func (self *MembersMetrica) GetValue() (float64, error) {
	memberCount, _ := self.Observable.GetStats()
	return float64(memberCount.Count()), nil
}

type FailuresMetrica struct {
	Observable Observable
}

func (self *FailuresMetrica) GetName() string {
	return "Failed to connect"
}
func (self *FailuresMetrica) GetUnits() string {
	return "Count"
}
func (self *FailuresMetrica) GetValue() (float64, error) {
	_, failures := self.Observable.GetStats()
	defer func() { failures.Clear() }()
	return float64(failures.Count()), nil
}
