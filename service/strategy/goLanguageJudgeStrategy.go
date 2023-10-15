/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 14:25:03
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-15 18:43:41
 * @FilePath: /xoj-backend/judge/strategy/impl/DefaultJudgeStrategy.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package strategy

import (
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"
	judgeinfomessageenum "github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/enums/JudgeInfoMessageEnum"
)

// Go 程序的判题策略
type GoLanguageJudgeStrategy struct {
}

// 执行判题
func (this GoLanguageJudgeStrategy) DoJudge(judgeContext JudgeContext) model.JudgeInfo {
	judgeInfoResponse := model.JudgeInfo{}
	executeCodeResponse := judgeContext.ExecuteCodeResponse //代码沙箱返回数据
	inputList := judgeContext.InputList
	judgeCaseList := judgeContext.JudgeCaseList
	quesionObj := judgeContext.Question

	outputList := executeCodeResponse.OutputList
	status := executeCodeResponse.Status   //执行状态
	message := executeCodeResponse.Message //对应status的信息描述
	judgeInfo := executeCodeResponse.JudgeInfo
	memory := judgeInfo.Memory //单位：byte
	time := judgeInfo.Time     //单位：ms

	// 实际消耗的内存和时间最大值
	judgeInfoResponse.Memory = memory
	judgeInfoResponse.Time = time

	// 判断沙箱执行状态是否正常
	if ok := CheckCodeSandboxResStatusOk(status, message, &judgeInfoResponse); !ok {
		return judgeInfoResponse
	}

	// // 判断沙箱执行状态是否正常
	// if !utils.CheckSame[int32]("判断沙箱执行的状态是否正常", executeCodeResponse.Status, codeexecstatusenum.SUCCEED.GetValue()) {
	// 	mylog.Log.Info("沙箱执行异常时返回的错误message: ", executeCodeResponse.Message)
	// 	// 编译错误
	// 	if executeCodeResponse.Status == codeexecstatusenum.COMPILE_FAIL.GetValue() || executeCodeResponse.Status == codeexecstatusenum.COMPILE_TIMEOUT_ERROR.GetValue() {
	// 		judgeInfoResponse.Message = judgeinfomessageenum.COMPILE_ERROR.GetValue()
	// 		return
	// 	}
	// 	// 运行错误
	// 	if executeCodeResponse.Status == codeexecstatusenum.RUN_FAIL.GetValue() || executeCodeResponse.Status == codeexecstatusenum.RUN_TIMEOUT_ERROR.GetValue() {
	// 		judgeInfoResponse.Message = judgeinfomessageenum.RUNTIME_ERROR.GetValue()
	// 		return
	// 	}
	// 	// 系统错误
	// 	if executeCodeResponse.Status == codeexecstatusenum.SYSTEM_ERROR.GetValue() {
	// 		judgeInfoResponse.Message = judgeinfomessageenum.SYSTEM_ERROR.GetValue()
	// 		return
	// 	}
	// }

	// 判断沙箱执行的结果输出数量是否和预期输出数量相等
	// if !utils.CheckSame[int]("判断沙箱执行的输入和输出数量是否一致", len(inputList), len(outputList)) {
	// 	// 答案错误
	// 	judgeInfoResponse.Message = judgeinfomessageenum.WRONG_ANSWER.GetValue()
	// 	return judgeInfoResponse
	// }
	if ok := CheckCodeSandboxResOutputLengthOk(inputList, outputList, &judgeInfoResponse); !ok {
		return judgeInfoResponse
	}

	// 依次判断每一项输出和预期输出是否相等
	// for i, len := 0, len(judgeCaseList); i < len; i++ {
	// 	judgeCase := judgeCaseList[i]
	// 	if !utils.CheckSame[string]("判断每项输出是否符合预期", judgeCase.Output, outputList[i]) {
	// 		// 答案错误
	// 		judgeInfoResponse.Message = judgeinfomessageenum.WRONG_ANSWER.GetValue()
	// 		return judgeInfoResponse
	// 	}
	// }
	if ok := CheckCodeSandboxResOutputRight(judgeCaseList, outputList, &judgeInfoResponse); !ok {
		return judgeInfoResponse
	}

	// 判断题目限制
	// judgeConfigStr := quesionObj.JudgeConfig
	// judgeConfig := question.JudgeConfig{}
	// if err := json.Unmarshal([]byte(judgeConfigStr), &judgeConfig); err != nil {
	// 	mylog.Log.Errorf("json.Unmarshal转换失败[%v]", judgeConfigStr)
	// 	// 系统错误
	// 	judgeInfoResponse.Message = judgeinfomessageenum.SYSTEM_ERROR.GetValue()
	// 	return judgeInfoResponse
	// }
	// needMemoryLimit := judgeConfig.MemoryLimit
	// needTimeLimit := judgeConfig.TimeLimit
	// if memory/1024 > needMemoryLimit {
	// 	mylog.Log.Errorf("实际使用内存=[%v]byte, 内存限制=[%v]KB", memory, needMemoryLimit)
	// 	// 内存溢出
	// 	judgeInfoResponse.Message = judgeinfomessageenum.MEMORY_LIMIT_EXCEEDED.GetValue()
	// 	return judgeInfoResponse
	// }
	// if time > needTimeLimit {
	// 	// 超时
	// 	judgeInfoResponse.Message = judgeinfomessageenum.TIME_LIMIT_EXCEEDED.GetValue()
	// 	return judgeInfoResponse
	// }
	if ok := CheckJudgeConfigPass(time, memory, quesionObj.JudgeConfig, &judgeInfoResponse); !ok {
		return judgeInfoResponse
	}

	// 判题后为正确
	judgeInfoResponse.Message = judgeinfomessageenum.ACCEPTED.GetValue()
	judgeInfoResponse.Detail = judgeInfoResponse.Message

	return judgeInfoResponse
}
