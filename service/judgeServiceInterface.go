/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 16:09:51
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-08 16:37:14
 */
package service

import "github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/entity"

// 判题服务
type JudgeServiceInterface interface {
	DoJudge(questionsubmitId int64) *entity.QuestionSubmit
}
