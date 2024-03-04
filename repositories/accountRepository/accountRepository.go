package accountrepository

import (
	"time"

	"pocketbase-demo/models"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
)

type accountRepository struct {
	PB *pocketbase.PocketBase
}

func NewAccountRepository(
	PB *pocketbase.PocketBase,
) accountRepository {
	return accountRepository{
		PB: PB,
	}
}

const column = `
	id, name, balance, created, updated
`

// GetAccountList gets a list of accounts
func (repo *accountRepository) GetAccountList(limit int, offset int) (accounts []models.Account, err error) {
	query := repo.PB.Dao().DB().NewQuery(
		"SELECT " + column + " FROM accounts LIMIT {:limit} OFFSET {:offset}",
	).Bind(dbx.Params{
		"limit":  limit,
		"offset": offset,
	})

	err = query.All(&accounts)
	return accounts, err
}

// IsAccountExistById checks if an account exists by its ID
func (repo *accountRepository) IsAccountExistById(accountId string) (account models.Account, exist bool) {
	query := repo.PB.Dao().DB().NewQuery(
		"SELECT " + column + " FROM accounts WHERE id = {:id}",
	).Bind(dbx.Params{
		"id": accountId,
	})

	err := query.One(&account)
	if err != nil {
		return account, false
	}

	return account, true
}

// IncrementDecrementBalance increments or decrements an account's balance
func (repo *accountRepository) IncrementDecrementBalance(accountId string, amount float64, increment bool, tx *daos.Dao) error {
	var query string
	if increment {
		query = `
			UPDATE accounts
			SET balance = balance + {:amount},
			updated = {:updated}
			WHERE id = {:id}
		`
	} else {
		query = `
			UPDATE accounts
			SET balance = balance - {:amount},
			updated = {:updated}
			WHERE id = {:id}
		`
	}

	param := dbx.Params{
		"id":      accountId,
		"amount":  amount,
		"updated": time.Now().Format("2006-01-02 15:04:05"),
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

// CreateAccount creates a new account
func (repo *accountRepository) CreateAccount(request models.Account) error {
	query := `
		INSERT INTO accounts (
			` + column + `
		) VALUES (
			{:id}, {:name}, {:balance}, {:created}, {:updated}
		)
	`

	param := dbx.Params{
		"id":      request.ID,
		"name":    request.Name,
		"balance": request.Balance,
		"created": time.Now().Format("2006-01-02 15:04:05"),
		"updated": time.Now().Format("2006-01-02 15:04:05"),
	}

	queryExecute := repo.PB.Dao().DB().NewQuery(query).Bind(param)
	_, err := queryExecute.Execute()
	return err
}
