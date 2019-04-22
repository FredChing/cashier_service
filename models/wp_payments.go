package models

import (
	"cashier_service/lib"
	"database/sql"
	"errors"
	"fmt"

	logs "github.com/cihub/seelog"
)

type WpPayments struct {
	Id             int64
	Mch_id         string
	Order_no       string
	Notifyurl      string
	Callbackurl    string
	Pay_memberid   int64
	Amount         float64
	Bankcode       string
	Orderid        string
	Out_trade_id   string
	Payment_status int
	Subject        string
	Created_at     int64
	Updated_at     int64
}

func (t *WpPayments) TableName() string {
	return "wp_payments"
}

func (payments *WpPayments) LoadByOrderId(order_id string) (err error) {
	select_str := fmt.Sprintf("SELECT `ID`,`mch_id`,`order_no`,`notifyurl`,`callbackurl`,`pay_memberid`,`amount`,`bankcode`,`orderid`,`out_trade_id`,`payment_status`,`subject`,`created_at`,`updated_at` FROM %s WHERE `orderid`='%s'", payments.TableName(), order_id)

	var (
		id             int64
		mch_id         string
		order_no       string
		notifyurl      string
		callbackurl    string
		pay_memberid   int64
		amount         float64
		bankcode       string
		orderid        string
		out_trade_id   string
		payment_status int
		subject        string
		created_at     int64
		updated_at     int64
	)

	err = lib.Mysql.QueryRow(select_str).Scan(
		&id, &mch_id, &order_no, &notifyurl,
		&callbackurl, &pay_memberid, &amount, &bankcode,
		&orderid, &out_trade_id, &payment_status,
		&subject, &created_at, &updated_at)
	switch {
	case err == sql.ErrNoRows:
		err = errors.New("not exist wp_payments")
		return
	case err != nil:
		_ = logs.Warnf("WpPayments::LoadByOrderId, select sql error, sql:%s, error:%s", select_str, err.Error())
		return
	default:
		payments.Id = id
		payments.Mch_id = mch_id
		payments.Order_no = order_no
		payments.Notifyurl = notifyurl
		payments.Callbackurl = callbackurl
		payments.Pay_memberid = pay_memberid
		payments.Amount = amount
		payments.Bankcode = bankcode
		payments.Orderid = orderid
		payments.Out_trade_id = out_trade_id
		payments.Payment_status = payment_status
		payments.Subject = subject
		payments.Created_at = created_at
		payments.Updated_at = updated_at
		return
	}
}
