package controllers

import (
	"cashier_service/models"
	"cashier_service/service"
	"encoding/base64"
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
	trade_id := this.GetString("trade_id")
	paymentService := &service.WpPaymentsService{}
	payment, err := paymentService.GetWpPaymentByTradeId(trade_id)
	if err != nil {
		this.Abort("订单号有误")
		return
	}
	logs.Infof("alipay_wap_ctrl::pay, get payment by trade_id success, trade_id:%s, payment:%v", trade_id, payment)
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
		p[key] = this.GetParams(f_key, payment)
	}

	r_s := []string{p["return_url"], p["source_return_url"]}
	p["return_url"] = strings.Join(r_s, "/")
	n_s := []string{p["notify_url"], p["source_notify_url"]}
	p["notify_url"] = strings.Join(n_s, "/")

	total_fee, _ := strconv.ParseFloat(p["total_amount"], 64)
	switch {
	case p["app_id"] == "", p["subject"] == "":
		this.Abort("pay_illegal_params")
		//this.OutputError(-1, errors.New("pay_illegal_params"))
		return
	case p["source_return_url"] == "":
		//this.OutputError(-1, errors.New("pay_illegal_params"))
		this.Abort("pay_illegal_params")
		return
	case p["source_notify_url"] == "":
		//this.OutputError(-1, errors.New("pay_illegal_params"))
		this.Abort("pay_illegal_params")
		return
	case total_fee <= 0:
		//this.OutputError(-1, errors.New("订单金额必须大于0"))
		this.Data["Content"] = "订单金额必须大于0"
		this.TplNames = "error/error.html"
		return
	}

	delete(p, "source_return_url")
	delete(p, "source_notify_url")

	pay := &service.AlipayService{
		Capture_account: strconv.Itoa(flag),
		Method:          service.ALIPAY_V2_METHOD_WAP,
		Sign_type:       service.ALIPAY_V2_SIGN_TYPE_RSA2,
	}
	_ = pay.PayInit(p)
	http_url, err := pay.GetUrl()
	logs.Infof("alipay::pay_new, pay:%v, http_url:%s",pay, http_url)
	if err != nil {
		//this.OutputError(-1, err)
		return
	}
	this.Redirect(http_url, 302)

	//this.OutputSuccess(http_url)
}

func (this *AlipayWap) GetParams(key string, payment *models.WpPayments) (str string) {
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
		str = payment.Subject
	case "return_url":
		domain := beego.AppConfig.String("domain")
		port := beego.AppConfig.String("httpport")
		str = "http://"+domain+":"+port+"/checkout/payback/alipay/return_url"
	case "notify_url":
		domain := beego.AppConfig.String("domain")
		port := beego.AppConfig.String("httpport")
		str = "http://"+domain+":"+port+"/checkout/payback/alipay/notify_url"
	case "timeout_express":
		str = "5m"
	case "total_amount":
		str = strconv.FormatFloat(payment.Amount, 'f', 2, 64)
	case "out_trade_no":
		str = payment.Out_trade_id
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
