package receptacle

type Receptacle struct {
	ID         int64  `json:"id" db:"id"`
	IDScenario int64  `json:"-" db:"id_scenario"`
	Name       string `json:"name" db:"name"`
	IDCard     int64  `json:"id_card" db:"id_card"`
}

// TODO model functions
