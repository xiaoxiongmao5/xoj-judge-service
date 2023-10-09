/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 14:25:03
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 17:42:30
 * @FilePath: /xoj-backend/judge/strategy/impl/DefaultJudgeStrategy.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package strategy

import (
	"encoding/json"

	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/dto/question"
	judgeinfomessageenum "github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/enums/JudgeInfoMessageEnum"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/utils"
)

// Go 程序的判题策略
type GoLanguageJudgeStrategy struct {
}

// 执行判题
func (this GoLanguageJudgeStrategy) DoJudge(judgeContext JudgeContext) (judgeInfoResponse model.JudgeInfo) {
	executeCodeResponse := judgeContext.ExecuteCodeResponse
	inputList := judgeContext.InputList
	judgeCaseList := judgeContext.JudgeCaseList
	quesionObj := judgeContext.Question

	outputList := executeCodeResponse.OutputList
	judgeInfo := executeCodeResponse.JudgeInfo
	memory := judgeInfo.Memory
	time := judgeInfo.Time

	// 实际消耗的内存和时间最大值
	judgeInfoResponse.Memory = memory
	judgeInfoResponse.Time = time

	// 判断沙箱执行状态是否正常
	if !utils.CheckSame[int32]("判断沙箱执行的状态是否正常", executeCodeResponse.Status, 1) {
		mylog.Log.Info("沙箱执行异常时返回的错误message: ", executeCodeResponse.Message)
		if executeCodeResponse.Status == 2 {
			// 编译错误
			judgeInfoResponse.Message = judgeinfomessageenum.COMPILE_ERROR.GetValue()
			return
		}
		if executeCodeResponse.Status == 3 {
			// 运行错误
			judgeInfoResponse.Message = judgeinfomessageenum.RUNTIME_ERROR.GetValue()
			return
		}
		if executeCodeResponse.Status == 4 {
			// 系统错误
			judgeInfoResponse.Message = judgeinfomessageenum.SYSTEM_ERROR.GetValue()
			return
		}
	}

	// 判断沙箱执行的结果输出数量是否和预期输出数量相等
	if !utils.CheckSame[int]("判断沙箱执行的输入和输出数量是否一致", len(inputList), len(outputList)) {
		// 答案错误
		judgeInfoResponse.Message = judgeinfomessageenum.WRONG_ANSWER.GetValue()
		return
	}

	// 依次判断每一项输出和预期输出是否相等
	for i, len := 0, len(judgeCaseList); i < len; i++ {
		judgeCase := judgeCaseList[i]
		if !utils.CheckSame[string]("判断每项输出是否符合预期", judgeCase.Output, outputList[i]) {
			// 答案错误
			judgeInfoResponse.Message = judgeinfomessageenum.WRONG_ANSWER.GetValue()
			return
		}
	}

	// 判断题目限制
	judgeConfigStr := quesionObj.JudgeConfig
	judgeConfig := question.JudgeConfig{}
	if err := json.Unmarshal([]byte(judgeConfigStr), &judgeConfig); err != nil {
		mylog.Log.Errorf("json.Unmarshal转换失败[%v]", judgeConfigStr)
		// 系统错误
		judgeInfoResponse.Message = judgeinfomessageenum.SYSTEM_ERROR.GetValue()
		return
	}
	needMemoryLimit := judgeConfig.MemoryLimit
	needTimeLimit := judgeConfig.TimeLimit
	if memory > needMemoryLimit {
		// 内存溢出
		judgeInfoResponse.Message = judgeinfomessageenum.MEMORY_LIMIT_EXCEEDED.GetValue()
		return
	}
	if time > needTimeLimit {
		// 超时
		judgeInfoResponse.Message = judgeinfomessageenum.TIME_LIMIT_EXCEEDED.GetValue()
		return
	}

	// 判题后为正确
	judgeInfoResponse.Message = judgeinfomessageenum.ACCEPTED.GetValue()
	return
}
