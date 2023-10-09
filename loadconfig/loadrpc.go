/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 20:04:28
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 10:49:55
 */
package loadconfig

import (
	"flag"
	"os"

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/rpc_api"
)

func init() {
	mylog.Log.Info("init begin: loadrpc")

	config.SetConsumerService(RpcQuestionClientImpl)
	config.SetConsumerService(RpcQuestionSubmitClientImpl)

	// 加载 Dubbo-go 的配置
	LoadDubboConfig()

	mylog.Log.Info("init end  : loadrpc")
}

var RpcQuestionClientImpl = new(rpc_api.QuestionClientImpl)
var RpcQuestionSubmitClientImpl = new(rpc_api.QuestionSubmitClientImpl)

// var replyQuestionGetByIdResp *rpc_api.QuestionGetByIdResp

// 设置环境变量
func SetOsEnv() {
	// 使用命令行参数来指定配置文件路径
	configFile := flag.String("config", "conf/dubbogo.yaml", "Path to Dubbo-go config file")
	flag.Parse()

	// 设置 DUBBO_GO_CONFIG_PATH 环境变量
	os.Setenv("DUBBO_GO_CONFIG_PATH", *configFile)
}

// 加载 Dubbo-go 的配置
func LoadDubboConfig() {
	SetOsEnv()
	// 加载 Dubbo-go 的配置文件，根据环境变量 DUBBO_GO_CONFIG_PATH 中指定的配置文件路径加载配置信息。配置文件通常包括 Dubbo 服务的注册中心地址、协议、端口等信息。
	if err := config.Load(); err != nil {
		panic(err)
	}
}
