package scenario

type Scenario struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	IDAuthor int64  `json:"-" db:"id_author"`
}

// TODO model functions
