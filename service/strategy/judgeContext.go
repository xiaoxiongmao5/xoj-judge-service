/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 14:02:29
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 17:28:46
 * @FilePath: /xoj-backend/judge/strategy/JudgeContext.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package strategy

import (
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/dto/question"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/entity"
)

type JudgeContext struct {
	ExecuteCodeResponse model.ExecuteCodeResponse
	InputList           []string
	JudgeCaseList       []question.JudgeCase
	Question            entity.Question
	QuestionSubmit      entity.QuestionSubmit
}
