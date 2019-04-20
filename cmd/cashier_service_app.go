package cmd

import (
	"cashier_service/service"
	"fmt"
	"net/http"
	"os"

	"github.com/astaxie/beego"
	logs "github.com/cihub/seelog"
)

var CashierServiceApp = cashierServiceApp{}

type cashierServiceApp struct {
	bStart bool
}

func (this *cashierServiceApp) Init() error {
	this.bStart = false
	currentPath, _ := os.Getwd()
	logConfigPath := beego.AppConfig.String("log_config_file")
	if logConfigPath != "" {
		logger, err := logs.LoggerFromConfigAsFile(currentPath + logConfigPath)

		if err != nil {
			fmt.Println("parse config.xml error", err)
		} else {
			logs.ReplaceLogger(logger)
		}
	}

	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}

	alipay_init()

	beego.BuildTemplate(beego.ViewsPath)

	return nil
}

func (this *cashierServiceApp) UnInit() error {
	this.bStart = false
	return nil
}

func (this *cashierServiceApp) Start() {
	this.bStart = true
	addr := beego.AppConfig.String("httpport")
	httpd := http.NewServeMux()
	httpd.Handle("/", beego.BeeApp.Handlers)

	listen := ":" + addr
	fmt.Println("server-started" + listen)
	err := http.ListenAndServe(listen, httpd)
	if err != nil {
		panic(err)
	}
}

func alipay_init() {
	service.AlipayV2RsaPath = beego.AppConfig.String("alipay::rsaPath")
	err := service.AlipayV2Init()
	if err != nil {
		panic(err)
	}
}
