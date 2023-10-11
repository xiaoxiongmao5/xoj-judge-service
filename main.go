/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 15:12:23
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-11 17:02:05
 * @FilePath: /xoj-judge-service/main.go
 */
package main

import (
	"sync"

	Octx "context"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/config"
	_ "github.com/xiaoxiongmao5/xoj/xoj-judge-service/config"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/consumer"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/middleware"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/myredis"
	_ "github.com/xiaoxiongmao5/xoj/xoj-judge-service/myrpc"
	_ "github.com/xiaoxiongmao5/xoj/xoj-judge-service/routers"
)

func init() {
	mylog.Log.Info("init begin: main")

	mylog.Log.Info("init end  : main")
}

func main() {
	defer mylog.Log.Writer().Close()
	defer myredis.Close(myredis.RedisCli)

	// 启动动态配置文件加载协程
	go config.LoadAppConfigDynamic()

	ctx := Octx.Background()

	// 创建互斥锁
	var mu sync.Mutex

	// 启动消费者协程
	go consumer.PopQuestionSubmit2Queue(ctx, myredis.RedisCli, &mu)

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

		// // 开启监控：Admin 管理后台
		// beego.BConfig.Listen.EnableAdmin = true
		// // 修改监听的地址和端口：
		// beego.BConfig.Listen.AdminAddr = "localhost"
		// beego.BConfig.Listen.AdminPort = 8089
	}

	// 全局异常捕获
	beego.BConfig.RecoverPanic = true
	beego.BConfig.RecoverFunc = middleware.ExceptionHandingMiddleware

	beego.Run()
}
