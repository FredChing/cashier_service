package controllers

import (
	"cashier_service/service"
	"encoding/base64"
	"errors"
	"github.com/astaxie/beego"
	logs "github.com/cihub/seelog"
	"strconv"
	"strings"
)

type AlipayWap struct {
	MyController
}

var flag = 1
func (this *AlipayWap) Pay() {
	//flag += 1
	//n := flag
	//if n == 5 {
	//	flag = 0
	//}
	//trade_id := this.GetString("trade_id")
	//return_url := this.GetString("return_url")
	//notify_url := this.GetString("notify_url")
	p := make(map[string]string)
	p_key := []string{
		"app_id",
		"out_trade_no",
		"total_amount",
		"notify_url",
		"return_url",
		"subject",
		"body",
		"timeout_express",
		"source_return_url",
		"source_notify_url",
	}

	for _, key := range p_key {
		var f_key string
		switch key {
		case "app_id":
			f_key = strconv.Itoa(flag) + "_" + key
		default:
			f_key = key
		}
		p[key] = this.GetParams(f_key)
	}

	r_s := []string{p["return_url"], p["source_return_url"]}
	p["return_url"] = strings.Join(r_s, "/")
	n_s := []string{p["notify_url"], p["source_notify_url"]}
	p["notify_url"] = strings.Join(n_s, "/")

	total_fee, _ := strconv.ParseFloat(p["total_amount"], 64)
	switch {
	case p["app_id"] == "", p["subject"] == "":
		this.OutputError(-1, errors.New("pay_illegal_params"))
		return
	case p["source_return_url"] == "":
		this.OutputError(-1, errors.New("pay_illegal_params"))
		return
	case p["source_notify_url"] == "":
		this.OutputError(-1, errors.New("pay_illegal_params"))
		return
	case total_fee <= 0:
		this.OutputError(-1, errors.New("订单金额必须大于0"))
		return
	}

	delete(p, "source_return_url")
	delete(p, "source_notify_url")
	p["out_trade_no"] = this.GetString("trade_id")

	pay := &service.AlipayService{
		Capture_account: strconv.Itoa(flag),
		Method:          service.ALIPAY_V2_METHOD_WAP,
		Sign_type:       service.ALIPAY_V2_SIGN_TYPE_RSA2,
	}
	_ = pay.PayInit(p)
	http_url, err := pay.GetUrl()
	logs.Infof("alipay::pay_new, pay:%v, http_url:%s",pay, http_url)
	if err != nil {
		this.OutputError(-1, err)
		return
	}

	this.OutputSuccess(http_url)
}

func (this *AlipayWap) GetParams(key string) (str string) {
	if strings.HasSuffix(key, "_app_id") {
		str = "pay::alipay_" + key
		str = beego.AppConfig.String(str)
		return
	}

	switch key {
	case "app_id":
		str = "pay::alipay_" + key
		str = beego.AppConfig.String(str)
	case "subject", "body":
		trade_id := this.GetString("trade_id")
		str = "测试订单" + trade_id
	case "return_url":
		str = "http://td.sandbox.wdwd.com/checkout/payback/alipay/return_url"
	case "notify_url":
		str = "http://td.sandbox.wdwd.com/checkout/payback/alipay/notify_url"
	case "timeout_express":
		str = "5m"
	case "total_amount":
		str = strconv.FormatFloat(0.01, 'f', 2, 64)
	case "out_trade_no":
		str = this.GetString("trade_id")
	case "source_return_url":
		str = base64.StdEncoding.EncodeToString([]byte(this.GetString("return_url")))
	case "source_notify_url":
		str = this.GetString("notify_url")
		if str == "" {
			str = this.GetString("return_url")
		}
		str = base64.StdEncoding.EncodeToString([]byte(str))
	}

	return
}
