package stat

type Stat struct {
	ID          int64  `json:"id" db:"id"`
	IDScenario  int64  `json:"-" db:"id_scenario"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	IDIcon      int64  `json:"id_icon" db:"id_icon"`
}

// TODO model functions
