package server

var state State

type State struct {
	Enabled *bool `json:"enabled,omitempty"`
	Volume  *int  `json:"volume,omitempty"`
}

/*func (s *State) Backup() {

}

func (s *State) Restore() {

}*/
