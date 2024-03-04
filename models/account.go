package models

type (
	Account struct {
		ID        string  `db:"id" json:"id"`
		Name      string  `db:"name" json:"name"`
		Balance   float64 `db:"balance" json:"balance"`
		CreatedAt string  `db:"created" json:"created"`
		UpdatedAt string  `db:"updated" json:"updated"`
	}

	RequestGetAccountList struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	RequestCreateAccount struct {
		Name string `json:"name" validate:"required"`
	}

	RequestGetAccountById struct {
		ID string `json:"id" validate:"required"`
	}
)
