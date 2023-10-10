/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 15:12:23
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-10 16:09:37
 * @FilePath: /xoj-judge-service/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"net/http"
	"sync"

	Octx "context"

	"github.com/beego/beego/v2/server/web/context"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/xiaoxiongmao5/xoj/xoj-judge-service/config"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/consumer"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/myredis"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/myresq"
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

	ctx := Octx.Background()

	// 创建互斥锁
	var mu sync.Mutex

	// 启动消费者协程
	consumer.PopQuestionSubmit2Queue(ctx, myredis.RedisCli, &mu)

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
	beego.BConfig.RecoverFunc = func(ctx *context.Context, config *beego.Config) {
		if err := recover(); err != nil {
			mylog.Log.Errorf("beego.BConfig.RecoverFunc err= %v \n", err)

			// 从 Context 中获取错误码和消息
			response, ok := ctx.Input.GetData("json").(*myresq.BaseResponse)
			if !ok {
				response = myresq.NewBaseResponse(500, "未知错误", nil)
			}

			// 将 JSON 响应写入 Context，并设置响应头
			ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
			ctx.Output.SetStatus(http.StatusOK)
			ctx.Output.JSON(response, false, false)
		}
	}

	beego.Run()
}
