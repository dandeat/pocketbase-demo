package repositories

import (
	"pocketbase-demo/models"

	"github.com/pocketbase/pocketbase/daos"
)

type AccountRepository interface {
	GetAccountList(limit int, offset int) (accounts []models.Account, err error)
	IsAccountExistById(accountId string) (account models.Account, exist bool)
	IncrementDecrementBalance(accountId string, amount float64, increment bool, tx *daos.Dao) error
	CreateAccount(request models.Account) error
}

type TransactionRepository interface {
	GetTransactionListByIndex(accountId, trxType string, limit int, offset int) (transactions []models.Transaction, err error)
	CreateTransaction(request models.Transaction, tx *daos.Dao) error
	GetTransactionByTransactionNumber(transactionNumber string) (transaction models.Transaction, exist bool)
}
