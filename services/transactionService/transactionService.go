package transactionservice

import (
	"log"
	"strconv"
	"time"

	"pocketbase-demo/models"
	"pocketbase-demo/services"
	"pocketbase-demo/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/daos"
)

type TransactionService struct {
	Service services.UseCaseService
}

func NewTransactionService(service services.UseCaseService) TransactionService {
	return TransactionService{
		Service: service,
	}
}

func (s TransactionService) GetTransactionList(ctx echo.Context) error {
	var (
		t       = time.Now()
		svcName = "[TransactionService]GetTransactionList"

		request      models.RequestGetTransactionList
		transactions []models.Transaction
	)

	// Validate JSON request
	if err := ctx.Bind(&request); err != nil {
		log.Println(svcName, "Error binding request:", err)
		response := utils.GenerateResponseJSON(t, "400", "Invalid request", nil)
		return ctx.JSON(400, response)
	}

	transactions, err := s.Service.TransactionRepo.GetTransactionListByIndex(
		request.AccountID, request.TransactionType, request.Limit, request.Offset)
	if err != nil {
		log.Println(svcName, "Error getting transaction list:", err)
		response := utils.GenerateResponseJSON(t, "500", "Internal server error", nil)
		return ctx.JSON(500, response)
	}

	response := utils.GenerateResponseJSON(t, "200", "Success", transactions)
	return ctx.JSON(200, response)
}

func (s TransactionService) GetTransactionByTransactionNumber(ctx echo.Context) error {
	var (
		t       = time.Now()
		svcName = "[TransactionService]GetTransactionByTransactionNumber"

		request     models.RequestGetTransactionByTrxNumber
		transaction models.Transaction
		exist       bool

		response models.ResponseJSON
	)

	// Validate JSON request
	if err := ctx.Bind(&request); err != nil {
		log.Println(svcName, "Error binding request:", err)
		response = utils.GenerateResponseJSON(t, "400", "Invalid request", nil)
		return ctx.JSON(400, response)
	}

	{ // Validation Layer
		// Validate TransactionNumber
		transaction, exist = s.Service.TransactionRepo.GetTransactionByTransactionNumber(request.TransactionNumber)
		if !exist {
			log.Println(svcName, "Transaction not found:", request.TransactionNumber)
			response = utils.GenerateResponseJSON(t, "404", "Transaction not found", nil)
			return ctx.JSON(404, response)
		}
	}

	response = utils.GenerateResponseJSON(t, "200", "Success", transaction)
	return ctx.JSON(200, response)
}

func (s TransactionService) TriggerDebit(ctx echo.Context) error {
	var (
		t       = time.Now()
		svcName = "[TransactionService]TriggerDebit"

		amountFloat float64
		err         error
		exist       bool

		account models.Account

		request  models.RequestTransactionDebit
		response models.ResponseJSON
	)

	// Validate JSON request
	if err := ctx.Bind(&request); err != nil {
		log.Println(svcName, "Error binding request:", err)
		response = utils.GenerateResponseJSON(t, "400", "Invalid request", nil)
		return ctx.JSON(400, response)
	}

	{ // Validation Layer
		// Validate Amount Format
		amountFloat, err = strconv.ParseFloat(request.Amount.Amount, 64)
		if err != nil {
			log.Println(svcName, "Error parsing amount:", err)
			response = utils.GenerateResponseJSON(t, "400", "Invalid amount", nil)
			return ctx.JSON(400, response)
		}

		// Validate Amount Value
		if amountFloat <= 0 {
			log.Println(svcName, "Invalid amount value:", request.Amount.Amount)
			response = utils.GenerateResponseJSON(t, "400", "Invalid amount value", nil)
			return ctx.JSON(400, response)
		}

		// Validate Currency
		if request.Amount.Currency != "USD" {
			log.Println(svcName, "Invalid currency:", request.Amount.Currency)
			response = utils.GenerateResponseJSON(t, "400", "Invalid currency", nil)
			return ctx.JSON(400, response)
		}

		// Validate AccountID
		account, exist = s.Service.AccountRepo.IsAccountExistById(request.AccountID)
		if !exist {
			log.Println(svcName, "Account not found:", request.AccountID)
			response = utils.GenerateResponseJSON(t, "404", "Account not found", nil)
			return ctx.JSON(404, response)
		}

		// Validate Balance
		if account.Balance < amountFloat {
			log.Println(svcName, "Insufficient balance:", account.Balance)
			response = utils.GenerateResponseJSON(t, "400", "Insufficient balance", nil)
			return ctx.JSON(400, response)
		}
	}

	err = s.Service.PB.Dao().RunInTransaction(func(tx *daos.Dao) error {

		uuidTrx := uuid.New().String()
		trxNumber := strconv.Itoa(int(t.Unix())) + uuidTrx

		err := s.Service.AccountRepo.IncrementDecrementBalance(request.AccountID, amountFloat, false, tx)
		if err != nil {
			log.Println(svcName, "Error decrementing balance:", err)
			return err
		}

		err = s.Service.TransactionRepo.CreateTransaction(models.Transaction{
			ID:                uuidTrx,
			TransactionNumber: trxNumber,
			TransactionType:   "DEBIT",
			AccountID:         request.AccountID,
			Amount:            amountFloat,
			Description:       request.Description,
			CreatedAt:         t.Format("2006-01-02 15:04:05"),
			UpdatedAt:         t.Format("2006-01-02 15:04:05"),
		}, tx)
		if err != nil {
			log.Println(svcName, "Error creating transaction:", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Println(svcName, "Error running transaction:", err)
		response = utils.GenerateResponseJSON(t, "500", "Internal server error", nil)
		return ctx.JSON(500, response)
	}

	response = utils.GenerateResponseJSON(t, "200", "Success Debit", nil)
	return ctx.JSON(200, response)
}

func (s TransactionService) TriggerCredit(ctx echo.Context) error {
	var (
		t       = time.Now()
		svcName = "[TransactionService]TriggerCredit"

		amountFloat float64
		err         error
		exist       bool

		request  models.RequestTransactionCredit
		response models.ResponseJSON
	)

	// Validate JSON request
	if err := ctx.Bind(&request); err != nil {
		log.Println(svcName, "Error binding request:", err)
		response = utils.GenerateResponseJSON(t, "400", "Invalid request", nil)
		return ctx.JSON(400, response)
	}

	{ // Validation Layer
		// Validate Amount Format
		amountFloat, err = strconv.ParseFloat(request.Amount.Amount, 64)
		if err != nil {
			log.Println(svcName, "Error parsing amount:", err)
			response = utils.GenerateResponseJSON(t, "400", "Invalid amount", nil)
			return ctx.JSON(400, response)
		}

		// Validate Amount Value
		if amountFloat <= 0 {
			log.Println(svcName, "Invalid amount value:", request.Amount.Amount)
			response = utils.GenerateResponseJSON(t, "400", "Invalid amount value", nil)
			return ctx.JSON(400, response)
		}

		// Validate Currency
		if request.Amount.Currency != "USD" {
			log.Println(svcName, "Invalid currency:", request.Amount.Currency)
			response = utils.GenerateResponseJSON(t, "400", "Invalid currency", nil)
			return ctx.JSON(400, response)
		}

		// Validate AccountID
		_, exist = s.Service.AccountRepo.IsAccountExistById(request.AccountID)
		if !exist {
			log.Println(svcName, "Account not found:", request.AccountID)
			response = utils.GenerateResponseJSON(t, "404", "Account not found", nil)
			return ctx.JSON(404, response)
		}
	}

	err = s.Service.PB.Dao().RunInTransaction(func(tx *daos.Dao) error {

		uuidTrx := uuid.New().String()
		trxNumber := strconv.Itoa(int(t.Unix())) + uuidTrx

		err := s.Service.AccountRepo.IncrementDecrementBalance(request.AccountID, amountFloat, true, tx)
		if err != nil {
			log.Println(svcName, "Error incrementing balance:", err)
			return err
		}

		err = s.Service.TransactionRepo.CreateTransaction(models.Transaction{
			ID:                uuidTrx,
			TransactionNumber: trxNumber,
			TransactionType:   "CREDIT",
			AccountID:         request.AccountID,
			Amount:            amountFloat,
			Description:       request.Description,
			CreatedAt:         t.Format("2006-01-02 15:04:05"),
			UpdatedAt:         t.Format("2006-01-02 15:04:05"),
		}, tx)
		if err != nil {
			log.Println(svcName, "Error creating transaction:", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Println(svcName, "Error running transaction:", err)
		response = utils.GenerateResponseJSON(t, "500", "Internal server error", nil)
		return ctx.JSON(500, response)
	}

	response = utils.GenerateResponseJSON(t, "200", "Success Credit", nil)
	return ctx.JSON(200, response)
}
