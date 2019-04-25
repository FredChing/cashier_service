package service

import (
	"cashier_service/lib"
	"cashier_service/models"
	logs "github.com/cihub/seelog"
)

type PayOrderService struct {
}

func (this *PayOrderService) GetPayOrderByTradeId(outTradeId string) (*models.PayOrder, error) {
	order := &models.PayOrder{}
	err := order.LoadByOrderId(outTradeId)
	if err != nil {
		_ = logs.Warnf("PayOrderService::GetPayOrderByTradeId, pay_order LoadByOrderId error, trade_id:%s, error:%s", outTradeId, err.Error())
		return nil, err
	}
	return order, nil
}

func (this *PayOrderService) AddPayOrder(tx *lib.CTx91dpays, pay_memberid string, pay_orderid string,pay_amount float64, pay_actualamount float64,
	subject string, out_trade_id string, notifyurl string, callbackurl string) (*models.PayOrder, error) {
	order := &models.PayOrder{}
	order.Pay_memberid = pay_memberid
	order.Pay_orderid = pay_orderid
	order.Pay_amount = pay_amount
	order.Pay_type = 1
	order.Pay_actualamount = pay_actualamount
	order.Pay_notifyurl = notifyurl
	order.Pay_callbackurl = callbackurl
	order.Pay_status = 0
	order.Pay_productname = subject
	order.Pay_tongdao = "alipay-wap"
	order.Pay_zh_tongdao = "支付宝-手机网站支付"
	order.Out_trade_id = out_trade_id
	order.Pay_ytongdao = "alipay-wap"
	order.Pay_yzh_tongdao = "支付宝-手机网站支付"
	order.Key = ""
	order.Account = ""
	order.Isdel = 0
	order.Ddlx = 0
	order.Pay_ytongdao = "alipay-wap"
	order.Pay_yzh_tongdao = "支付宝-手机网站支付"
	order.Xx = 0
	order.Attach = ""
	order.Pay_channel_account = ""
	err := order.Insert(tx)
	if err != nil {
		_ = logs.Warnf("PayOrderService::AddPayOrder, pay_order Insert error, pay_order:%v, error:%s", order, err.Error())
		return nil, err
	}
	return order, nil
}

func (this *PayOrderService) UpdateOrderPayStatusPending(out_trade_id string) error {
	order := &models.PayOrder{}
	err := order.UpdateOrderPayStatusPending(out_trade_id)
	if err != nil {
		_ = logs.Warnf("PayOrderService::UpdateOrderPayStatusPending, pay_order update error, out_trade_id:%s, error:%s", out_trade_id, err.Error())
		return err
	}
	return nil
}

func (this *PayOrderService) UpdateOrderPayStatusSuccess(out_trade_id string) error {
	order := &models.PayOrder{}
	err := order.UpdateOrderPayStatusSuccess(out_trade_id)
	if err != nil {
		_ = logs.Warnf("PayOrderService::UpdateOrderPayStatusSuccess, pay_order update error, out_trade_id:%s, error:%s", out_trade_id, err.Error())
		return err
	}
	return nil
}
