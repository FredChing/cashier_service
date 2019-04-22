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
