package models

type AccountModel struct {
	Id        int    `json:"id,omitempty" db:"id"`
	FirstName string `json:"firstName"  db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
}
