package location

type Location struct {
	ID         int64  `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	IDScenario int64  `json:"-" db:"id_scenario"`
	IDCardA    *int64 `json:"id_card_a,omitempty" db:"id_card_a"`
	IDCardB    *int64 `json:"id_card_b,omitempty" db:"id_card_b"`
	IDCardC    *int64 `json:"id_card_c,omitempty" db:"id_card_c"`
	IDCardD    *int64 `json:"id_card_d,omitempty" db:"id_card_d"`
	IDCardE    *int64 `json:"id_card_e,omitempty" db:"id_card_e"`
	IDCardF    *int64 `json:"id_card_f,omitempty" db:"id_card_f"`
	IDCardG    *int64 `json:"id_card_g,omitempty" db:"id_card_g"`
	IDCardH    *int64 `json:"id_card_h,omitempty" db:"id_card_h"`
}

// TODO model functions
