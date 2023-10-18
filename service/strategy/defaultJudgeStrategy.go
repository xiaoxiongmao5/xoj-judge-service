/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 14:25:03
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-17 23:09:53
 * @FilePath: /xoj-backend/judge/strategy/impl/DefaultJudgeStrategy.go
 * @Description: 默认判题策略
 */
package strategy

import (
	"encoding/json"
	"fmt"

	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/dto/question"
	codeexecstatusenum "github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/enums/CodeExecStatusEnum"
	judgeinfomessageenum "github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/enums/JudgeInfoMessageEnum"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/utils"
)

// 默认判题策略
type DefaultJudgeStrategy struct {
}

// 判断沙箱执行状态是否正常
func CheckCodeSandboxResStatusOk(status int32, message string, judgeInfoResponse *model.JudgeInfo) bool {
	// status := executeCodeResponse.Status   //执行状态
	// message := executeCodeResponse.Message //对应status的信息描述

	if !utils.CheckSame[int32]("判断沙箱返回数据-执行状态[status]是否正常", status, codeexecstatusenum.SUCCEED.GetValue()) {
		judgeInfoResponse.Detail = message
		switch status {
		case codeexecstatusenum.COMPILE_FAIL.GetValue(): //编译错误
			judgeInfoResponse.Message = judgeinfomessageenum.COMPILE_ERROR.GetValue()
			return false
		case codeexecstatusenum.COMPILE_TIMEOUT_ERROR.GetValue(): //编译超时
			judgeInfoResponse.Message = judgeinfomessageenum.COMPILE_TIME_LIMIT_EXCEEDED.GetValue()
			return false
		case codeexecstatusenum.RUN_FAIL.GetValue(): // 运行错误
			judgeInfoResponse.Message = judgeinfomessageenum.RUNTIME_ERROR.GetValue()
			return false
		case codeexecstatusenum.RUN_TIMEOUT_ERROR.GetValue(): //运行超时
			judgeInfoResponse.Message = judgeinfomessageenum.RUN_TIME_LIMIT_EXCEEDED.GetValue()
			return false
		case codeexecstatusenum.OUT_OF_MEMORY_ERROR.GetValue(): //内存不足
			judgeInfoResponse.Message = judgeinfomessageenum.OUT_OF_MEMORY.GetValue()
			return false
		case codeexecstatusenum.SYSTEM_ERROR.GetValue(): //系统错误
			judgeInfoResponse.Message = judgeinfomessageenum.SYSTEM_ERROR.GetValue()
			return false
		}
	}

	return true
}

// 判断沙箱执行的结果输出数量是否和预期输出数量相等
func CheckCodeSandboxResOutputLengthOk(inputList, outputList []string, judgeInfoResponse *model.JudgeInfo) bool {
	// 先判断沙箱执行的结果输出数量是否和预期输出数量相等
	if !utils.CheckSame[int]("判断沙箱执行的输入和输出数量是否一致", len(inputList), len(outputList)) {
		// 答案错误
		judgeInfoResponse.Message = judgeinfomessageenum.WRONG_ANSWER.GetValue()
		judgeInfoResponse.Detail = "执行的输出条目与输入用例数量不等"
		return false
	}
	return true
}

// 判断每一项输出和预期输出是否相等
func CheckCodeSandboxResOutputRight(judgeCaseList []question.JudgeCase, outputList []string, judgeInfoResponse *model.JudgeInfo) bool {
	for i, len := 0, len(judgeCaseList); i < len; i++ {
		judgeCase := judgeCaseList[i]
		if !utils.CheckSame[string]("判断每项输出是否符合预期", judgeCase.Output, outputList[i]) {
			// 答案错误
			judgeInfoResponse.Message = judgeinfomessageenum.WRONG_ANSWER.GetValue()
			detail := fmt.Sprintf("执行输出结果错误，输入=[%s]时,预期结果=[%s],你的结果=[%s]", judgeCase.Input, judgeCase.Output, outputList[i])
			judgeInfoResponse.Detail = detail
			return false
		}
	}
	return true
}

// 判断是否通过题目限制要求（内存、耗时）
func CheckJudgeConfigPass(time, memory int64, judgeConfigStr string, judgeInfoResponse *model.JudgeInfo) bool {
	judgeConfig := question.JudgeConfig{}
	if err := json.Unmarshal([]byte(judgeConfigStr), &judgeConfig); err != nil {
		mylog.Log.Errorf("json.Unmarshal转换失败[%v]", judgeConfigStr)
		// 系统错误
		judgeInfoResponse.Message = judgeinfomessageenum.SYSTEM_ERROR.GetValue()
		judgeInfoResponse.Detail = judgeInfoResponse.Message
		return false
	}
	needMemoryLimit := judgeConfig.MemoryLimit
	needTimeLimit := judgeConfig.TimeLimit
	if memory/1024 > needMemoryLimit {
		detail := fmt.Sprintf("实际使用内存=[%v]byte, 内存限制=[%v]KB", memory, needMemoryLimit)
		mylog.Log.Errorf(detail)
		// 内存溢出
		judgeInfoResponse.Message = judgeinfomessageenum.MEMORY_LIMIT_EXCEEDED.GetValue()
		judgeInfoResponse.Detail = detail
		return false
	}
	if time > needTimeLimit {
		detail := fmt.Sprintf("实际使用时间=[%v]ms, 时间限制=[%v]ms", time, needTimeLimit)
		// 超题限时
		judgeInfoResponse.Message = judgeinfomessageenum.TIME_LIMIT_EXCEEDED.GetValue()
		judgeInfoResponse.Detail = detail
		return false
	}
	return true
}

// 执行判题
func (this DefaultJudgeStrategy) DoJudge(judgeContext JudgeContext) model.JudgeInfo {
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

	// 先判断沙箱执行的结果输出数量是否和预期输出数量相等
	if ok := CheckCodeSandboxResOutputLengthOk(inputList, outputList, &judgeInfoResponse); !ok {
		return judgeInfoResponse
	}

	// 依次判断每一项输出和预期输出是否相等
	if ok := CheckCodeSandboxResOutputRight(judgeCaseList, outputList, &judgeInfoResponse); !ok {
		return judgeInfoResponse
	}

	// 判断题目限制
	if ok := CheckJudgeConfigPass(time, memory, quesionObj.JudgeConfig, &judgeInfoResponse); !ok {
		return judgeInfoResponse
	}

	// 判题后为正确
	judgeInfoResponse.Message = judgeinfomessageenum.ACCEPTED.GetValue()
	judgeInfoResponse.Detail = judgeInfoResponse.Message

	return judgeInfoResponse
}
