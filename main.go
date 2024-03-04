package main

import (
	"log"
	"os"

	accountrepository "pocketbase-demo/repositories/accountRepository"
	transactionrepository "pocketbase-demo/repositories/transactionRepository"
	"pocketbase-demo/services"
	accountservice "pocketbase-demo/services/accountService"
	transactionservice "pocketbase-demo/services/transactionService"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()
	// err := collections.InitCollection(app)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// declare Repo
	accountRepo := accountrepository.NewAccountRepository(app)
	transactionRepo := transactionrepository.NewTransactionRepository(app)
	useCaseSvc := services.NewUseCaseService(
		app,
		&accountRepo,
		&transactionRepo,
	)

	//

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))

		// declare API
		apiGroup := e.Router.Group("/api")

		// declare API for account
		accountAPI := apiGroup.Group("/account")
		accountSvc := accountservice.NewAccountService(useCaseSvc)
		accountAPI.POST("/list", accountSvc.GetAccountList)
		accountAPI.POST("/add", accountSvc.CreateAccount)

		// declare API for transaction
		transactionAPI := apiGroup.Group("/transaction")
		transactionSvc := transactionservice.NewTransactionService(useCaseSvc)
		transactionAPI.POST("/list", transactionSvc.GetTransactionList)
		transactionAPI.POST("/trigger-debit", transactionSvc.TriggerDebit)
		transactionAPI.POST("/trigger-credit", transactionSvc.TriggerCredit)

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
