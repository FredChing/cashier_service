package service

import (
	"cashier_service/models"
	logs "github.com/cihub/seelog"
)

type WpPaymentsService struct {
}

//CheckTradeByTradeId 根据订单号获取订单信息
func (this *WpPaymentsService) GetWpPaymentByTradeId(trade_id string) (*models.WpPayments, error) {
	payment := &models.WpPayments{}
	err := payment.LoadByOrderId(trade_id)
	if err != nil {
		_ = logs.Warnf("WpPaymentsService::GetWpPaymentByTradeID, payment LoadByOrderId error, trade_id:%s, error:%s", trade_id, err.Error())
		return nil, err
	}
	return payment, nil
}

func (this *WpPaymentsService) AddPayment(mch_id string, order_no string,orderid string, out_trade_id string,
	subject string, amount float64, notifyurl string, callbackurl string) (*models.WpPayments, error) {
	payment := &models.WpPayments{}
	payment.Mch_id = mch_id
	payment.Order_no = order_no
	payment.Notifyurl = notifyurl
	payment.Callbackurl = callbackurl
	payment.Pay_memberid = 0
	payment.Amount = amount
	payment.Bankcode = ""
	payment.Orderid = orderid
	payment.Out_trade_id = out_trade_id
	payment.Payment_status = 0
	payment.Subject = subject
	err := payment.Insert()
	if err != nil {
		_ = logs.Warnf("WpPaymentsService::AddPayment, payment Insert error, payment:%v, error:%s", payment, err.Error())
		return nil, err
	}
	return payment, nil
}

func (this *WpPaymentsService) UpdatePaymentStatusSuccess(order_no string) error {
	payment := &models.WpPayments{}
	err := payment.UpdatePaymentStatusSuccess(order_no)
	if err != nil {
		_ = logs.Warnf("WpPaymentsService::UpdatePaymentStatusSuccess, payment update error, orderid:%s, error:%s", order_no, err.Error())
		return err
	}
	return nil
}
