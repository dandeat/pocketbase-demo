package accountservice

import (
	"log"
	"time"

	"pocketbase-demo/models"
	"pocketbase-demo/services"
	"pocketbase-demo/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type AccountService struct {
	Service services.UseCaseService
}

func NewAccountService(
	service services.UseCaseService,
) AccountService {
	return AccountService{
		Service: service,
	}
}

// GetAccountList is a function to get account list
func (s AccountService) GetAccountList(ctx echo.Context) error {
	var (
		svcName = "[AccountService]GetAccountList"
		t       = time.Now()

		request  models.RequestGetAccountList
		response models.ResponseJSON
		accounts []models.Account
	)

	// Validate JSON request
	if err := ctx.Bind(&request); err != nil {
		log.Println(svcName, "Error binding request:", err)
		response = utils.GenerateResponseJSON(t, "400", "Invalid request", nil)
		return ctx.JSON(400, response)
	}

	accounts, err := s.Service.AccountRepo.GetAccountList(request.Limit, request.Offset)
	if err != nil {
		log.Println(svcName, "Error getting account list:", err)
		response = utils.GenerateResponseJSON(t, "500", "Internal server error", nil)
		return ctx.JSON(500, response)
	}

	response = utils.GenerateResponseJSON(t, "200", "Success", accounts)
	return ctx.JSON(200, response)
}

// GetAccountByAccountId is a function to get account by account id
func (s AccountService) GetAccountByAccountId(ctx echo.Context) error {
	var (
		svcName = "[AccountService]GetAccountByAccountId"
		t       = time.Now()

		request models.RequestGetAccountById
		account models.Account
		exist   bool

		response models.ResponseJSON
	)

	// Validate JSON request
	if err := ctx.Bind(&request); err != nil {
		log.Println(svcName, "Error binding request:", err)
		response = utils.GenerateResponseJSON(t, "400", "Invalid request", nil)
		return ctx.JSON(400, response)
	}

	account, exist = s.Service.AccountRepo.IsAccountExistById(request.ID)
	if !exist {
		response = utils.GenerateResponseJSON(t, "404", "Account not found", nil)
		return ctx.JSON(404, response)
	}

	response = utils.GenerateResponseJSON(t, "200", "Success", account)
	return ctx.JSON(200, response)
}

// CreateAccount is a function to create account
func (s AccountService) CreateAccount(ctx echo.Context) error {
	var (
		svcName = "[AccountService]CreateAccount"
		t       = time.Now()

		request  models.RequestCreateAccount
		response models.ResponseJSON
	)

	// Validate JSON request
	if err := ctx.Bind(&request); err != nil {
		log.Println(svcName, "Error binding request:", err)
		response = utils.GenerateResponseJSON(t, "400", "Invalid request", nil)
		return ctx.JSON(400, response)
	}

	// Create Account
	err := s.Service.AccountRepo.CreateAccount(models.Account{
		ID:        uuid.New().String(),
		Name:      request.Name,
		Balance:   0,
		CreatedAt: t.Format("2006-01-02 15:04:05"),
		UpdatedAt: t.Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		log.Println(svcName, "Error creating account:", err)
		response = utils.GenerateResponseJSON(t, "500", "Internal server error", nil)
		return ctx.JSON(500, response)
	}

	response = utils.GenerateResponseJSON(t, "200", "Success", nil)
	return ctx.JSON(200, response)
}
