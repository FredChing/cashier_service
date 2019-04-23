package service

import (
	"cashier_service/lib"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	logs "github.com/cihub/seelog"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
	"time"
)

type AlipayService struct {
	PubParams       map[string]string
	Method          string
	Sign_type       string
	Capture_account string
	param_str       string
}

const (
	ALIPAY_V2_GATEWAY = "https://openapi.alipay.com/gateway.do"
	ALIPAY_V2_FORMAT  = "JSON"
	ALIPAY_V2_CHARSET = "utf-8"
	ALIPAY_V2_VERSION = "1.0"

	ALIPAY_V2_METHOD_PAGE    = "alipay.trade.page.pay"
	ALIPAY_V2_METHOD_WAP  = "alipay.trade.wap.pay"
	ALIPAY_V2_SIGN_TYPE_RSA2 = "RSA2"

	ALIPAY_V2_TRADE_WAIT_PAY = "WAIT_BUYER_PAY"
	ALIPAY_V2_TRADE_CLOSED   = "TRADE_CLOSED"
	ALIPAY_V2_TRADE_SUCC     = "TRADE_SUCCESS"
)

var (
	alipayV2Method2Pcode_M = map[string]string{
		ALIPAY_V2_METHOD_WAP:  "QUICK_WAP_WAY",
		ALIPAY_V2_METHOD_PAGE: "FAST_INSTANT_TRADE_PAY",
	}

	alipayV2BizParam_pay_M = map[string]bool{
		"out_trade_no":         true,
		"body":                 true,
		"subject":              true,
		"total_amount":         true,
		"product_code":         true,
		"timeout_express":      false,
		"goods_type":           false,
		"passback_params":      false,
		"extend_params":        false,
		"enable_pay_channels":  false,
		"disable_pay_channels": false,
	}

	alipayV2PubParam_M = map[string]bool{
		"app_id":      true,
		"method":      true,
		"format":      false,
		"charset":     true,
		"sign_type":   true,
		"timestamp":   true,
		"version":     true,
		"notify_url":  false,
		"return_url":  false,
		"biz_content": true,
	}

	alipayV2BizParam_MM = map[string](map[string]bool){
		ALIPAY_V2_METHOD_WAP:              alipayV2BizParam_pay_M,
		ALIPAY_V2_METHOD_PAGE:             alipayV2BizParam_pay_M,
	}

	AlipayV2RsaPath      string
	alipayV2PublicKey_M  map[string]*rsa.PublicKey
	alipayV2PrivateKey_M map[string]*rsa.PrivateKey
)

func AlipayV2Init() (err error) {
	capture_acc_arr := []string{"1"}
	alipayV2PrivateKey_M = make(map[string]*rsa.PrivateKey)
	alipayV2PublicKey_M = make(map[string]*rsa.PublicKey)

	for index, _ := range capture_acc_arr {
		capture_acc := capture_acc_arr[index]

		// 解析私钥
		privateKeyPath := AlipayV2RsaPath + capture_acc + "_rsa_private_key.pem"
		publicKeyPath := AlipayV2RsaPath + capture_acc + "_rsa_public_key.pem"
		logs.Infof("alipay_service::AlipayV2Init, privateKeyPath:%s, publicKeyPath:%s", privateKeyPath, publicKeyPath)
		kt, _ := ioutil.ReadFile(privateKeyPath)
		block, _ := pem.Decode(kt)
		alipayV2PrivateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

		// 解析公钥
		public_kt, _ := ioutil.ReadFile(publicKeyPath)
		PEMBlock, _ := pem.Decode([]byte(public_kt))
		if PEMBlock == nil {
			_ = logs.Warnf("alipay_service::AlipayV2Init, pem.Decode failed, privateKeyPath:%s, publicKeyPath:%s", privateKeyPath, publicKeyPath)
			return errors.New("Could not parse Public Key PEM")
		}
		if PEMBlock.Type != "PUBLIC KEY" {
			_ = logs.Warnf("alipay_service::AlipayV2Init, Found wrong public key type, privateKeyPath:%s, publicKeyPath:%s", privateKeyPath, publicKeyPath)
			return errors.New("wrong public key")
		}

		pubkey, err := x509.ParsePKIXPublicKey(PEMBlock.Bytes)
		if err != nil {
			_ = logs.Warnf("alipay_service::AlipayV2Init, x509.ParsePKIXPublicKey failed, privateKeyPath:%s, publicKeyPath:%s, error:%s", privateKeyPath, publicKeyPath, err.Error())
			return err
		}
		alipayV2PublicKey := pubkey.(*rsa.PublicKey)

		alipayV2PrivateKey_M[capture_acc] = alipayV2PrivateKey
		alipayV2PublicKey_M[capture_acc] = alipayV2PublicKey
		logs.Infof("aaaaaaaaaaaa::AlipayV2Init, alipayV2PrivateKey:%v, alipayV2PrivateKey_M:%v", alipayV2PrivateKey, alipayV2PrivateKey_M)
		logs.Infof("bbbbbbbbbbbb::AlipayV2Init, alipayV2PublicKey:%v, alipayV2PublicKey_M:%v", alipayV2PublicKey, alipayV2PublicKey_M)
	}

	return
}

//删除value为空的,key全转为小写
func filterParams(params map[string]string) map[string]string {
	p := make(map[string]string, 0)
	for k, v := range params {
		if v != "" {
			p[strings.ToLower(k)] = v
		}
	}
	return p
}

func (alp2 *AlipayService) PayInit(params map[string]string) (err error) {
	params = filterParams(params)
	pub_params, biz_params := alp2.split_params(params)

	now := time.Now()
	pub_params["timestamp"] = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
	pub_params["method"] = alp2.Method
	pub_params["version"] = ALIPAY_V2_VERSION
	pub_params["format"] = ALIPAY_V2_FORMAT
	pub_params["charset"] = ALIPAY_V2_CHARSET
	pub_params["sign_type"] = ALIPAY_V2_SIGN_TYPE_RSA2
	biz_params["product_code"] = alipayV2Method2Pcode_M[alp2.Method]

	err = alp2.check_params(biz_params, alipayV2BizParam_pay_M)
	if err != nil {
		_ = logs.Warnf("alipay_service::PayInit, check_params biz_params alipayV2BizParam_pay_M error, biz_params:%v, alipayV2BizParam_pay_M:%v, error:%s", biz_params, biz_params, err.Error())
		return
	}
	biz_params = filterParams(biz_params)
	byteBiz, err := json.Marshal(biz_params)
	if err != nil {
		_ = logs.Warnf("alipay_service::PayInit, json.Marshal biz_params error, biz_params:%v, alipayV2BizParam_pay_M:%v, error:%s", biz_params, biz_params, err.Error())
		return
	}
	pub_params["biz_content"] = string(byteBiz)

	err = alp2.check_params(pub_params, alipayV2PubParam_M)
	if err != nil {
		_ = logs.Warnf("alipay_service::PayInit, check_params pub_params alipayV2PubParam_M error, pub_params:%v, alipayV2PubParam_M:%v, error:%s", pub_params, alipayV2PubParam_M, err.Error())
		return
	}
	pub_params = filterParams(pub_params)
	alp2.PubParams = pub_params

	sign, err := alp2.sign(alp2.PubParams)
	if err != nil {
		_ = logs.Warnf("alipay_service::PayInit, alp2.sign pub_params error, pub_params:%v, error:%s", alp2.PubParams, err.Error())
		return
	}
	alp2.PubParams["sign"] = sign
	alp2.param_str = alp2.urlencode(alp2.PubParams)

	return
}

func (alp2 *AlipayService) split_params(params map[string]string) (pub_params, biz_params map[string]string) {
	pub_params = make(map[string]string)
	biz_params = make(map[string]string)
	for key, value := range params {
		if _, ok := alipayV2PubParam_M[key]; ok {
			pub_params[key] = value
		}

		if _, ok := alipayV2BizParam_pay_M[key]; ok {
			biz_params[key] = value
		}
	}
	return
}

func (alp2 *AlipayService) check_params(params map[string]string, param_must_M map[string]bool) (err error) {
	for k, is_must := range param_must_M {
		if _, ok := params[k]; !ok {
			if is_must {
				err = errors.New("no arg " + k)
				return
			}
		}
	}
	return
}

func (alp2 *AlipayService) sign(params map[string]string) (sign string, err error) {
	var keys = make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		value := strings.TrimSpace(params[key])
		if key != "sign" && value != "" {
			pList = append(pList, key+"="+value)
		}
	}
	src := strings.Join(pList, "&")
	println(src)
	var byteSign []byte
	byteSign, err = rsaEncryptPrivateKey(src, alp2.get_private_key())
	if err != nil {
		return
	}
	sign = base64.StdEncoding.EncodeToString(byteSign)

	return
}

func (alp2 *AlipayService) get_private_key() (rsa_private_key *rsa.PrivateKey) {
	if _rsa_private_key, ok := alipayV2PrivateKey_M[alp2.Capture_account]; ok {
		rsa_private_key = _rsa_private_key
	} else {
		rsa_private_key = alipayV2PrivateKey_M["1"]
	}

	return
}

func (alp2 *AlipayService) get_public_key() (rsa_public_key *rsa.PublicKey) {
	if _rsa_public_key, ok := alipayV2PublicKey_M[alp2.Capture_account]; ok {
		rsa_public_key = _rsa_public_key
	} else {
		rsa_public_key = alipayV2PublicKey_M["1"]
	}
	return
}

func rsaEncryptPrivateKey(str string, privateKey *rsa.PrivateKey) ([]byte, error) {
	h := sha256.New()
	h.Write([]byte(str))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, h.Sum(nil))
	if err != nil {
		_ = logs.Warnf("alipay_service::rsaEncryptPrivateKey, rsa.SignPKCS1v15 error, str:%s, error:%s", str, err.Error())
	}

	return signature, err
}

func (alp2 *AlipayService) urlencode(params map[string]string) string {
	values := url.Values{}
	for key, value := range params {
		if value != "" {
			values.Set(key, value)
		}
	}
	return values.Encode()
}

func (alp2 *AlipayService) GetUrl() (info string, err error) {
	if alp2.param_str != "" {
		return ALIPAY_V2_GATEWAY + "?" + alp2.param_str, nil
	}

	sign, err := alp2.sign(alp2.PubParams)
	if err != nil {
		return
	}
	alp2.PubParams["sign"] = sign

	return ALIPAY_V2_GATEWAY + "?" + alp2.urlencode(alp2.PubParams), nil
}

func (alp2 *AlipayService) Sign(secret string, params map[string]string) string {
	var keys = make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		value := strings.TrimSpace(params[key])
		if key != "sign" && value != "" {
			pList = append(pList, key+"="+value)
		}
	}
	src := strings.Join(pList, "&") + "&" + secret
	logs.Infof("alipay_service::Sign, sign befor, src:%s", src)
	sign := lib.MD5(src)
	sign = strings.ToUpper(sign)
	logs.Infof("alipay_service::Sign, sign befor, src:%s，after sign:%s", src, sign)
	return sign
}

func (alp2 *AlipayService) VerifySign(values url.Values) (verified bool, err error) {
	return alp2.verify_sign(values)
}

func (alp2 *AlipayService) verify_sign(values url.Values) (verified bool, err error) {
	byteSign, err := base64.StdEncoding.DecodeString(values.Get("sign"))
	if err != nil {
		return
	}
	sign := string(byteSign)

	sign_type := values.Get("sign_type")
	if sign == "" || sign_type == "" {
		err = errors.New("no arg sign or sign_type")
		return
	}
	alp2.Sign_type = sign_type

	keys := make([]string, 0)
	for key, value := range values {
		if key == "sign" || key == "sign_type" {
			continue
		}

		if len(value) > 0 {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	pLst := make([]string, 0)
	for _, key := range keys {
		value := strings.TrimSpace(values.Get(key))
		if value != "" {
			pLst = append(pLst, key+"="+value)
		}
	}
	str := strings.Join(pLst, "&")

	err = rsa2Verify(str, alp2.get_public_key(), sign)
	if err != nil {
		return
	}
	verified = true

	return
}

func rsa2Verify(str string, pub *rsa.PublicKey, signature string) (err error) {
	h := sha256.New()
	h.Write([]byte(str))

	hashed := h.Sum(nil)
	sign_byte := []byte(signature)

	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], sign_byte)
	return
}

func (this *AlipayService) RequestNotifyUrl(url string, params *map[string]string) error {
	req := httplib.Post(url)
	//设置超时时间：链接3s超时，读写60s超时
	connectTimeout := time.Duration(3 * time.Second)
	readWriteTimeout := time.Duration(60 * time.Second)
	req.SetTimeout(connectTimeout, readWriteTimeout)
	var s []string
	for k, v := range *params {
		req.Param(k, v)
		s = append(s, k+":"+v+", ")
	}
	paramStr := strings.Join(s, "")
	paramStr = paramStr[:len(paramStr)-2]
	byteData, err := req.Bytes()
	if err != nil {
		_ = logs.Warnf("AlipayService::RequestNotifyUrl error, url:%s, params:%s, response:%s", url, paramStr, string(byteData))
		return err
	}
	//result := &WeiXuanResult{}
	//err = json.Unmarshal(byteData, &result)
	//if err != nil {
	//	logs.Warnf("WeixuanService::RequestWeiXuanApi, result data json unmarshal fail:api_url:%s, params:%s, error:%s", api_url, paramStr, err.Error())
	//	return nil, err
	//}
	logs.Infof("AlipayService::RequestNotifyUrl success, url:%s, params:%s, response:%+v", url, paramStr, string(byteData))
	return err
}
