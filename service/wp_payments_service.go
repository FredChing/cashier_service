package service

import (
	"cashier_service/lib"
	"cashier_service/models"
	logs "github.com/cihub/seelog"
)

type WpPaymentsService struct {
}

//CheckTradeByTradeId 根据订单号获取订单信息
func (this *WpPaymentsService) GetWpPaymentByTradeId(outTradeId string) (*models.WpPayments, error) {
	payment := &models.WpPayments{}
	err := payment.LoadByOrderId(outTradeId)
	if err != nil {
		_ = logs.Warnf("WpPaymentsService::GetWpPaymentByTradeID, payment LoadByOrderId error, outTradeId:%s, error:%s", outTradeId, err.Error())
		return nil, err
	}
	return payment, nil
}

func (this *WpPaymentsService) AddPayment(tx *lib.CTx, mch_id string, order_no string, orderid string, out_trade_id string,
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
	err := payment.Insert(tx)
	if err != nil {
		_ = logs.Warnf("WpPaymentsService::AddPayment, payment Insert error, payment:%v, error:%s", payment, err.Error())
		return nil, err
	}
	return payment, nil
}

func (this *WpPaymentsService) UpdatePaymentStatusSuccess(out_trade_id string) error {
	payment := &models.WpPayments{}
	err := payment.UpdatePaymentStatusSuccess(out_trade_id)
	if err != nil {
		_ = logs.Warnf("WpPaymentsService::UpdatePaymentStatusSuccess, payment update error, out_trade_id:%s, error:%s", out_trade_id, err.Error())
		return err
	}
	return nil
}
