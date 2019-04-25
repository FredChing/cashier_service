package models

import (
	"cashier_service/lib"
	"database/sql"
	"errors"
	"fmt"
	"time"

	logs "github.com/cihub/seelog"
)

type PayOrder struct {
	Id                  int64
	Pay_memberid        string
	Pay_orderid         string
	Qrcode_id           int64
	Pay_type            int
	Pay_amount          float64
	Pay_poundage        float64
	Pay_actualamount    float64
	Pay_applydate       int64
	Pay_successdate     int64
	Pay_bankcode        string
	Pay_notifyurl       string
	Pay_callbackurl     string
	Pay_bankname        string
	Pay_status          int
	Pay_productname     string
	Pay_tongdao         string
	Pay_zh_tongdao      string
	Pay_tjurl           string
	Out_trade_id        string
	Num                 int
	Memberid            string
	Key                 string
	Account             string
	Isdel               int
	Ddlx                int64
	Pay_ytongdao        string
	Pay_yzh_tongdao     string
	Xx                  int
	Attach              string
	Pay_channel_account string
	Cost                float64
	Cost_rate           float64
	Account_id          int64
	Channel_id          int64
	T                   int
	Last_reissue_time   int64
	Lock_status         int
}

func (t *PayOrder) TableName() string {
	return "pay_order"
}

func (order *PayOrder) LoadByOrderId(outTradeId string) (err error) {
	select_str := fmt.Sprintf("SELECT `pay_memberid`,`pay_orderid`,`pay_type`,`pay_amount`,`pay_actualamount`,`pay_applydate`,`pay_successdate`,`out_trade_id`,`pay_status` FROM %s WHERE `out_trade_id`='%s'", order.TableName(), outTradeId)

	var (
		pay_memberid     string
		pay_orderid      string
		pay_type         int
		pay_amount       float64
		pay_actualamount float64
		pay_applydate    int64
		pay_successdate  int64
		out_trade_id     string
		pay_status       int
	)

	err = lib.Mysql_91dpays.QueryRow(select_str).Scan(
		&pay_memberid, &pay_orderid, &pay_type, &pay_amount,
		&pay_actualamount, &pay_applydate, &pay_successdate, &out_trade_id,
		&pay_status)
	switch {
	case err == sql.ErrNoRows:
		err = errors.New("not exist pay_order")
		return
	case err != nil:
		_ = logs.Warnf("PayOrder::LoadByOrderId, select sql error, sql:%s, error:%s", select_str, err.Error())
		return
	default:
		order.Pay_memberid = pay_memberid
		order.Pay_orderid = pay_orderid
		order.Pay_type = pay_type
		order.Pay_amount = pay_amount
		order.Pay_actualamount = pay_actualamount
		order.Pay_applydate = pay_applydate
		order.Pay_successdate = pay_successdate
		order.Out_trade_id = out_trade_id
		order.Pay_status = pay_status
		return
	}
}

//Insert 添加数据
func (this *PayOrder) Insert(tx *lib.CTx91dpays) error {
	if this.Pay_applydate == 0 {
		this.Pay_applydate = time.Now().Unix()
	}
	sql_str := fmt.Sprintf("INSERT INTO `pay_order` (`pay_memberid`, `pay_orderid`, `qrcode_id`, `pay_type`, `pay_amount`, `pay_poundage`, `pay_actualamount`, `pay_applydate`, `pay_successdate`, `pay_bankcode`, `pay_notifyurl`, `pay_callbackurl`, `pay_bankname`, `pay_status`, `pay_productname`, `pay_tongdao`, `pay_zh_tongdao`, `pay_tjurl`, `out_trade_id`, `num`, `memberid`, `key`, `account`, `isdel`, `ddlx`, `pay_ytongdao`, `pay_yzh_tongdao`, `xx`, `attach`, `pay_channel_account`, `cost`, `cost_rate`, `account_id`, `channel_id`, `t`, `last_reissue_time`, `lock_status`) VALUES ('%s', '%s', %d, %d, %g, '0.0000', %g, %d, 0, '', '%s', '%s', '', '0', '%s', '%s', '%s', '', '%s', '0', '1', '', '', '0', '0', '%s', '%s', '0', '', '', '0.0000', '0.0000', '35', '1', '0', '11', '0')",
		this.Pay_memberid, this.Pay_orderid, this.Qrcode_id, this.Pay_type, this.Pay_amount, this.Pay_actualamount, this.Pay_applydate, this.Pay_notifyurl, this.Pay_callbackurl, this.Pay_productname, this.Pay_tongdao, this.Pay_zh_tongdao, this.Out_trade_id, this.Pay_ytongdao, this.Pay_yzh_tongdao)
	_, err := tx.Exec(sql_str)
	if err != nil {
		tx.Rollback()
		_ = logs.Warnf("PayOrder::Insert, insert failed, sql:%s, error:%s", sql_str, err.Error())
		return err
	}
	return nil
}

func (this *PayOrder) UpdateOrderPayStatusPending(out_trade_id string) error {
	sql_str := fmt.Sprintf("UPDATE %s SET `pay_status`=1 WHERE out_trade_id='%s'",
		this.TableName(), out_trade_id)
	tx, err := lib.NewTx91dpays()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sql_str)
	if err != nil {
		tx.Rollback()
		_ = logs.Warnf("PayOrder::UpdateOrderPayStatusPending, update failed , sql:%s, error:%s", sql_str, err.Error())
		return err
	}
	_ = tx.Commit()
	return nil
}

func (this *PayOrder) UpdateOrderPayStatusSuccess(out_trade_id string) error {
	sql_str := fmt.Sprintf("UPDATE %s SET `pay_status`=2, `pay_successdate`=%d WHERE out_trade_id='%s'",
		this.TableName(), time.Now().Unix(), out_trade_id)
	tx, err := lib.NewTx91dpays()
	if err != nil {
		return err
	}
	_, err = tx.Exec(sql_str)
	if err != nil {
		tx.Rollback()
		_ = logs.Warnf("PayOrder::UpdateOrderPayStatusSuccess, update failed , sql:%s, error:%s", sql_str, err.Error())
		return err
	}
	_ = tx.Commit()
	return nil
}
