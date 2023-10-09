/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 14:23:36
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 13:12:00
 */
package strategy

import "github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"

// 判题策略
type JudgeStrategyInterface interface {
	DoJudge(judgeContext JudgeContext) model.JudgeInfo
}
