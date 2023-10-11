/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 15:34:11
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 13:43:06
 */
package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/myresq"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/service"
)

type MainController struct {
	beego.Controller
}

//	@Summary		根据 id 判题
//	@Description	根据 id 判题
//	@Tags			判题
//	@Accept			application/x-www-form-urlencoded
//	@Produce		application/json
//	@Param			id	query		int		true	"id"
//	@Success		200	{object}	object	"响应数据"
//	@Router			/dojudge [get]
func (this MainController) DoJudge() {
	id, err := this.GetInt64("id")
	if err != nil || id <= 0 {
		myresq.Abort(this.Ctx, myresq.PARAMS_ERROR, "")
		return
	}

	questionSubmitObj := service.DoJudge(this.Ctx, id)
	myresq.Success(this.Ctx, questionSubmitObj)
}
