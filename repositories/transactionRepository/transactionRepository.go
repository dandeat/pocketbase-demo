package transactionrepository

import (
	"pocketbase-demo/models"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
)

type transactionRepository struct {
	PB *pocketbase.PocketBase
}

func NewTransactionRepository(
	PB *pocketbase.PocketBase,
) transactionRepository {
	return transactionRepository{
		PB: PB,
	}
}

const column = `
	id, transaction_number, transaction_type, account_id, 
	amount, description, created, updated
`

// GetTransactionListByIndex gets a list of transactions by index
func (repo *transactionRepository) GetTransactionListByIndex(accountId, trxType string, limit int, offset int) (transactions []models.Transaction, err error) {

	query := `
		SELECT ` + column + ` FROM transactions 
		WHERE 1=1
	`
	param := dbx.Params{}

	if accountId != "" {
		query += " AND account_id = {:account_id}"
		param["account_id"] = accountId
	}

	if trxType != "" {
		query += " AND transaction_type = {:transaction_type}"
		param["transaction_type"] = trxType
	}

	query += " LIMIT {:limit} OFFSET {:offset}"
	param["limit"] = limit
	param["offset"] = offset

	err = repo.PB.Dao().DB().NewQuery(query).Bind(param).All(&transactions)
	return transactions, err
}

// CreateTransaction creates a new transaction
func (repo *transactionRepository) CreateTransaction(request models.Transaction, tx *daos.Dao) error {

	query := `
		INSERT INTO transactions (
			` + column + `
		) VALUES (
			{:id}, {:transaction_number}, {:transaction_type}, {:account_id},
			{:amount}, {:description}, {:created}, {:updated}
		)
	`

	param := dbx.Params{
		"id":                 request.ID,
		"transaction_number": request.TransactionNumber,
		"transaction_type":   request.TransactionType,
		"account_id":         request.AccountID,
		"amount":             request.Amount,
		"description":        request.Description,
		"created":            request.CreatedAt,
		"updated":            request.UpdatedAt,
	}

	var queryExecute *dbx.Query
	if tx == nil {
		queryExecute = repo.PB.Dao().DB().NewQuery(query).Bind(param)
	} else {
		queryExecute = tx.DB().NewQuery(query).Bind(param)
	}

	_, err := queryExecute.Execute()
	return err
}

// GetTransactionByTransactionNumber gets a transaction by its transaction number
func (repo *transactionRepository) GetTransactionByTransactionNumber(transactionNumber string) (transaction models.Transaction, exist bool) {
	query := repo.PB.Dao().DB().NewQuery(
		"SELECT " + column + " FROM transactions WHERE transaction_number = {:transaction_number}",
	).Bind(dbx.Params{
		"transaction_number": transactionNumber,
	})

	err := query.One(&transaction)
	if err != nil {
		return transaction, false
	}

	return transaction, true
}
