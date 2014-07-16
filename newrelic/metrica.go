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
	normal, _ := self.Observable.GetStats()
	return float64(normal), nil
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
	return float64(failures), nil
}
