package controllers

import (
	"cashier_service/models"
	"cashier_service/service"
	"encoding/base64"
	"github.com/astaxie/beego"
	logs "github.com/cihub/seelog"
	"strconv"
	"strings"
	"time"
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
	mch_id := this.GetString("mch_id")
	trade_id := this.GetString("trade_id")
	return_url := this.GetString("return_url")
	notify_url := this.GetString("notify_url")
	amount := this.GetString("amount")
	product_name := this.GetString("product_name")
	sign := this.GetString("ac_sign")
	logs.Infof("alipay_wap_ctrl::Pay, recv data, requestBody:%s, mch_id:%s, trade_id:%s, return_url:%s, notify_url:%s, "+
		"amount:%s, product_name:%s, csign:%s", string(this.Ctx.Input.RequestBody), mch_id, trade_id, return_url, notify_url, amount, product_name, sign)
	if len(trade_id) <= 0 || len(return_url) <= 0 || len(mch_id) <= 0 || len(amount) <= 0 || len(product_name) <= 0 { //参数校验
		this.Abort("缺少参数")
		return
	}

	sign_key := beego.AppConfig.String("pay::alipay_" + strconv.Itoa(flag) + "_key") //签名key
	logs.Infof("alipay_service::Pay, sign_key:%s, key:%s", sign_key, "pay::alipay_"+strconv.Itoa(flag)+"_key")
	params := map[string]string{
		"mch_id":       mch_id,
		"trade_id":     trade_id,
		"amount":       amount,
		"product_name": product_name,
		"return_url":   return_url,
		"notify_url":   notify_url,
	}
	//签名校验
	alipayService := &service.AlipayService{}
	sign_str := alipayService.Sign(sign_key, params)
	if sign != sign_str {
		this.Abort("签名错误")
		return
	}

	//add
	paymentService := &service.WpPaymentsService{}
	order_no := string(time.Now().Format("20060102150405")) + strconv.Itoa(time.Now().Nanosecond())
	amount_f, _ := strconv.ParseFloat(amount, 64)
	payment, err := paymentService.AddPayment(mch_id, order_no, trade_id, order_no, product_name, amount_f, notify_url, return_url)
	if err != nil {
		this.Abort("出错啦")
		return
	}
	logs.Infof("alipay_wap_ctrl::pay, AddPayment success, trade_id:%s", trade_id)
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
	logs.Infof("alipay::pay_new, pay:%v, http_url:%s", pay, http_url)
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
		str = "http://" + domain + ":" + port + "/checkout/payback/alipay/return_url"
	case "notify_url":
		domain := beego.AppConfig.String("domain")
		port := beego.AppConfig.String("httpport")
		str = "http://" + domain + ":" + port + "/checkout/payback/alipay/notify_url"
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
