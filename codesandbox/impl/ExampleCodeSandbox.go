/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 12:29:56
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 17:59:18
 */
package impl

import (
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"
	judgeinfomessageenum "github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/enums/JudgeInfoMessageEnum"
	questionsubmitstatusenum "github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/enums/QuestionSubmitStatusEnum"
)

type ExampleCodeSandbox struct {
}

func (this ExampleCodeSandbox) ExecuteCode(executeCodeRequest model.ExecuteCodeRequest) (model.ExecuteCodeResponse, error) {
	return model.ExecuteCodeResponse{
		OutputList: executeCodeRequest.InputList,
		Message:    "测试执行成功",
		Status:     questionsubmitstatusenum.SUCCEED.GetValue(),
		JudgeInfo: model.JudgeInfo{
			Message: judgeinfomessageenum.ACCEPTED.GetText(),
			Memory:  100,
			Time:    100,
		},
	}, nil
}
