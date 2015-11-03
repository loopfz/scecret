package element

type Element struct {
	ID          int64  `json:"-" db:"id"`
	IDScenario  int64  `json:"-" db:"id_scenario"`
	Number      int    `json:"number" db:"number"`
	Description string `json:"description" db:"description"`
	IDCard      int64  `json:"id_card" db:"id_card"`
}

// TODO model functions
