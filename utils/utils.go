package utils

import (
	"time"

	"pocketbase-demo/models"
)

func GenerateResponseJSON(t time.Time, status string, message string, data interface{}) models.ResponseJSON {
	return models.ResponseJSON{
		ResponseCode:     status,
		ResponseDesc:     message,
		ResponseDateTime: t.Format("2006-01-02 15:04:05"),
		Result:           data,
	}
}
