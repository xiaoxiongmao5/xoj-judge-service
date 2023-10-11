/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 15:12:23
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 13:45:01
 */
//	@APIVersion			1.0.0
//	@Title				beego Test API
//	@Description		beego has a very cool tools to autogenerate documents for your API
//	@Contact			astaxie@gmail.com
//	@TermsOfServiceUrl	http://beego.me/
//	@License			Apache 2.0
//	@LicenseUrl			http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/controllers"
)

func init() {
	beego.CtrlGet("/dojudge", controllers.MainController.DoJudge)
}
