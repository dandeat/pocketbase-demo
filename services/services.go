package services

import (
	"pocketbase-demo/repositories"

	"github.com/pocketbase/pocketbase"
)

type UseCaseService struct {
	PB              *pocketbase.PocketBase
	AccountRepo     repositories.AccountRepository
	TransactionRepo repositories.TransactionRepository
}

func NewUseCaseService(
	PB *pocketbase.PocketBase,
	accountRepo repositories.AccountRepository,
	transactionRepo repositories.TransactionRepository,
) UseCaseService {
	return UseCaseService{
		PB:              PB,
		AccountRepo:     accountRepo,
		TransactionRepo: transactionRepo,
	}
}
