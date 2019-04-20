// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/astaxie/beego"
	logs "github.com/cihub/seelog"
)

var CmdServer = &Command{
	UsageLine: "server [start | stop | status | reload]",
	Short:     "controller the pay server",
	Long: `
    start              start the server.
    stop               stop the server.
    status             show the server status.
    reload             reload server by send HUP signal.
`,
}

func init() {
	CmdServer.Run = server_controller
	beego.CopyRequestBody = true
	os.MkdirAll("log", os.ModePerm)
}

func server_controller(cmd *Command, args []string) {

	// if len(args) != 1 {
	// 	fmt.Println("[ERRO] action?")
	// 	os.Exit(2)
	// }

	switch args[0] {
	case "start":
		start_server()
	case "start-console":
		start_server()
	case "stop":
	case "status":
	case "reload":
	}

	fmt.Println("[SUCC] New application successfully created!", args)
	os.Exit(0)
}

func start_server() {
	wg := &sync.WaitGroup{}
	//	wg.Add(1)
	StartCustomerApp(wg)
	wg.Wait()

}

func StartCustomerApp(wg *sync.WaitGroup) {
	logs.Info("service start_server...")

	defer func() {
		wg.Done()
		logs.Info("service close")
	}()

	_ = CashierServiceApp.Init()
	CashierServiceApp.Start()
	logs.Info("service start_server, end...")
}
