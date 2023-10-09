/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 13:58:53
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 15:04:14
 */
package service

import (
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/service/strategy"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/utils"
)

// 判题管理（简化调用）
type JudgeManager struct {
}

// 选择对应的判题策略-执行判题
func (this JudgeManager) DoJudge(judgeContext strategy.JudgeContext) model.JudgeInfo {
	// 设置默认的判题策略
	judgeStrategy := strategy.DefaultJudgeStrategy{}

	language := judgeContext.QuestionSubmit.Language

	if utils.CheckSame[string]("判断是否为Go语言", language, "go") {
		judgeStrategy := strategy.GoLanguageJudgeStrategy{}
		return judgeStrategy.DoJudge(judgeContext)
	}

	return judgeStrategy.DoJudge(judgeContext)
}
