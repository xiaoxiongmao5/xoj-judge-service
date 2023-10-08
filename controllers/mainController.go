/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 15:34:11
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-08 15:34:36
 */
package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (this MainController) DoJudge() {
	fmt.Println("hi xj")
}
