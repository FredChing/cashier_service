package routers

import (
	"cashier_service/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/checkout/pay/alipay", &controllers.Alipay{}, "post:Pay") //支付宝页面扫码支付
	beego.Router("/checkout/pay/alipay_wap", &controllers.AlipayWap{}, "post:Pay") //支付宝手机网站支付
	beego.Router("/checkout/payback/alipay/return_url/:callback", &controllers.AlipayBack{}, "*:HandleReturn")
	beego.Router("/checkout/payback/alipay/notify_url/:callback", &controllers.AlipayBack{}, "*:HandleNotify")
}
