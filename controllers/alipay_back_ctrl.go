package controllers

import (
	"cashier_service/service"
	"encoding/base64"
	"github.com/astaxie/beego"
	"net/url"
	"strconv"
	"strings"

	logs "github.com/cihub/seelog"
)

type AlipayBack struct {
	MyController
}

func (this *AlipayBack) HandleReturn() {
	out_trade_no := this.GetString("out_trade_no")
	logs.Infof("AlipayBack::HandleReturn, out_trade_no:%s", out_trade_no)

	uri := &url.Values{}
	callback := this.GetString(":callback")
	callback = decodeToUrl(callback)

	//uri.Set("trade_id", trade_id)
	logs.Infof("AlipayBack::HandleReturn, out_trade_no:%s, return_url:%s", out_trade_no, callback+"?"+uri.Encode())
	this.Redirect(callback+"?"+uri.Encode(), 302)
}

func decodeToUrl(base64_str string) string {
	data, err := base64.StdEncoding.DecodeString(base64_str)
	if err != nil {
		return ""
	}
	return string(data)
}

func (this *AlipayBack) HandleNotify() {
	logs.Info("AlipayBack::HandleNotify")
	var vals url.Values
	if this.Ctx.Input.Request.Form == nil {
		this.Ctx.Input.Request.ParseForm()
	}
	vals = this.Ctx.Input.Request.Form

	out_trade_no := vals.Get("out_trade_no")
	logs.Infof("AlipayBack::HandleNotify, out_trade_no:%s", out_trade_no)

	alipayService := &service.AlipayService{
		Capture_account: "1",
	}
	verified, err := alipayService.VerifySign(vals)
	if err != nil {
		_ = logs.Warnf("AlipayBack::HandleNotify, VerifySign failed, vals:%v, error:%s", vals, err)
		return
	}
	logs.Infof("AlipayBack::HandleNotify, VerifySign success, vals:%v, verify is %v", vals.Encode(), verified)

	this.handleNotify(vals)
}

func (this *AlipayBack) handleNotify(vals url.Values) {
	logs.Infof("AlipayBack::HandleNotify, notify success, vals:%v", vals)

	switch vals.Get("trade_status") {
	case service.ALIPAY_V2_TRADE_SUCC:
		callback := this.GetString(":callback")
		callback = decodeToUrl(callback)
		logs.Infof("AlipayBack::HandleNotify, trade success, callbackUrl:%s, vals:%v", callback, vals)
		params := map[string]string{
			"bizOrderNo": vals.Get("out_trade_no"),
			"orderNo":    vals.Get("trade_no"),
			"amount":     vals.Get("receipt_amount"),
			"payTime":    vals.Get("gmt_payment"),
			"status":     "0",
		}
		sign_key := beego.AppConfig.String("pay::alipay_" + strconv.Itoa(flag) + "_key") //签名key
		alipayService := &service.AlipayService{}
		sign_str := alipayService.Sign(sign_key, params)
		params["sign"] = sign_str
		err := alipayService.RequestNotifyUrl(callback, &params)
		if err != nil {
			_ = logs.Warnf("AlipayBack::HandleNotify, RequestNotifyUrl failed, callbackUrl:%s, params:%v, vals:%v", callback, params, vals)
		}
		this.Ctx.WriteString("success")
		this.view("success")
	case service.ALIPAY_V2_TRADE_WAIT_PAY, service.ALIPAY_V2_TRADE_CLOSED:
		this.view("success")
	default:
		this.view("failed")
	}
}

func (this *AlipayBack) CheckInWechat() (isTrue bool) {

	ua := this.Ctx.Request.Header["User-Agent"]
	if len(ua) > 0 {
		isTrue = strings.Contains(ua[0], "MicroMessenger")
	} else {
		isTrue = false
	}
	return
}

func (this *AlipayBack) view(result string) {
	this.Ctx.Output.Status = 200
	this.Ctx.Output.Body([]byte(strings.ToLower(result)))
}
