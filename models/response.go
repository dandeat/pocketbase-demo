package models

type ResponseJSON struct {
	ResponseCode     string      `json:"response_code"`
	ResponseDesc     string      `json:"response_desc"`
	ResponseDateTime string      `json:"response_date_time"`
	Result           interface{} `json:"result"`
}
