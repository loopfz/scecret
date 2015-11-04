package models

type Element struct {
	ID          int64  `json:"-" db:"id"`
	IDScenario  int64  `json:"-" db:"id_scenario"`
	Number      int    `json:"number" db:"number"`
	Description string `json:"description" db:"description"`
	Notes       string `json:"notes" db:"notes"`
	IDCard      int64  `json:"id_card" db:"id_card"`
}

// TODO model functions
