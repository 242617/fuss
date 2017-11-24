package server

var state *State

type State struct {
	Enabled bool `json:"enabled"`
	Volume  int  `json:"volume"`
	Left    int  `json:"left"`
	Right   int  `json:"right"`
}

func NewState() *State {
	return &State{
		Enabled: false,
		Volume:  100,
	}
}

/*func (s *State) Backup() {

}

func (s *State) Restore() {

}*/
