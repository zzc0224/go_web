package controller

import (
	"bluebell_backend/models"
)

type _ResponseLogin struct {
	Code    MyCode                   `json:"code"`    //业务响应状态码
	Message string                   `json:"message"` //提示信息
	Data    []*models.ApiLoginDetail `json:"data"`    //数据
}
