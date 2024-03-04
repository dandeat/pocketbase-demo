package models

type (
	Transaction struct {
		ID                string  `db:"id" json:"id"`
		TransactionNumber string  `db:"transaction_number" json:"transaction_number"`
		TransactionType   string  `db:"transaction_type" json:"transaction_type"` // DEBIT, CREDIT
		AccountID         string  `db:"account_id" json:"account_id"`
		Amount            float64 `db:"amount" json:"amount"`
		Description       string  `db:"description" json:"description"`
		CreatedAt         string  `db:"created" json:"created"`
		UpdatedAt         string  `db:"updated" json:"updated"`
	}

	DataAmount struct {
		Amount   string `json:"amount" validate:"required"`   // 00.00
		Currency string `json:"currency" validate:"required"` // USD, EUR, IDR
	}

	RequestTransactionDebit struct {
		AccountID   string     `json:"account_id" validate:"required"`
		Amount      DataAmount `json:"amount" validate:"required,dive"`
		Description string     `json:"description"`
	}

	RequestTransactionCredit struct {
		AccountID   string     `json:"account_id" validate:"required"`
		Amount      DataAmount `json:"amount" validate:"required,dive"`
		Description string     `json:"description"`
	}

	RequestGetTransactionByTrxNumber struct {
		TransactionNumber string `json:"transaction_number" validate:"required"`
	}

	RequestGetTransactionList struct {
		AccountID       string `json:"account_id"`
		TransactionType string `json:"transaction_type"`

		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}
)
