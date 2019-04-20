package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
)

// 代码库操作
type MyController struct {
	beego.Controller
}

//通用返回结果
type APIResult struct {
	Ecode  int         `json:"ecode"` //0表示成功，其他错误各个接口定义 在非0的情况下 Emsg 要有原因
	Emsg   string      `json:"emsg"`
	Result interface{} `json:"result"`
}

func (this *MyController) Prepare() {
}

func (m *APIResult) MarshalJSON() ([]byte, error) {
	v := map[string]interface{}{
		"ecode":  m.Ecode,
		"emsg":   m.Emsg,
		"result": m.Result,
	}
	return json.Marshal(v)
}

func (this *MyController) OutputError(code int, err error) {
	this.Data["json"] = &APIResult{Ecode: code, Emsg: err.Error(), Result: nil}
	this.ServeJson()
}

func (this *MyController) OutputSuccess(result interface{}) {
	this.Data["json"] = &APIResult{Ecode: 0, Emsg: "success", Result: result}
	this.ServeJson()
}