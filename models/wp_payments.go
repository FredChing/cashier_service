package models

import (
	"cashier_service/lib"
	"database/sql"
	"errors"
	"fmt"
	"time"

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
	select_str := fmt.Sprintf("SELECT `ID`,`mch_id`,`order_no`,`notifyurl`,`callbackurl`,`pay_memberid`,`amount`,`bankcode`,`orderid`,`out_trade_id`,`payment_status`,`subject` FROM %s WHERE `orderid`='%s'", payments.TableName(), order_id)

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
	)

	err = lib.Mysql.QueryRow(select_str).Scan(
		&id, &mch_id, &order_no, &notifyurl,
		&callbackurl, &pay_memberid, &amount, &bankcode,
		&orderid, &out_trade_id, &payment_status,
		&subject)
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
		return
	}
}

//Insert 添加数据
func (this *WpPayments) Insert() error {
	if this.Created_at == 0 {
		this.Created_at = time.Now().Unix()
		this.Updated_at = this.Created_at
	}
	tx, err := lib.NewTx()
	if err != nil {
		return err
	}
	_, err = tx.Insert(this)
	if err != nil {
		tx.Rollback()
		_ = logs.Warnf("WpPayments::Insert, insert failed, error:%s", err.Error())
		return err
	}
	tx.Commit()
	return nil
}

//UpdateCallInterface 更新调用提现接口为已发送状态
func (this *WpPayments) UpdatePaymentStatusSuccess(orderid string) error {
	sql_str := fmt.Sprintf("UPDATE %s SET `payment_status`=1, `updated_at`=%d WHERE orderid='%s'",
		this.TableName(), time.Now().Unix(), orderid)
	tx, err := lib.NewTx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sql_str)
	if err != nil {
		tx.Rollback()
		_ = logs.Warnf("WpPayments::UpdatePaymentStatusSuccess, update failed , sql:%s, error:%s", sql_str, err.Error())
		return err
	}
	_ = tx.Commit()
	return nil
}
