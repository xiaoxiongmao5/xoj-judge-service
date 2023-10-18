/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 13:27:42
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-18 00:42:39
 */
package service

import (
	"context"
	"encoding/json"
	"errors"

	beeContext "github.com/beego/beego/v2/server/web/context"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/dto/question"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/entity"
	judgeinfomessageenum "github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/enums/JudgeInfoMessageEnum"
	questionsubmitlanguageenum "github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/enums/QuestionSubmitLanguageEnum"
	questionsubmitstatusenum "github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/enums/QuestionSubmitStatusEnum"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/myerror"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/myresq"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/myrpc"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/rpc_api"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/service/strategy"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/utils"
)

func UpdateQuestionSubmitObj(ctx context.Context, beeCtx *beeContext.Context, questionSubmitObj entity.QuestionSubmit, rpcQuestionsubmitObj *rpc_api.RpcQuestionSubmitObj) {
	rpcQuestionsubmitObj.Status = questionSubmitObj.Status
	rpcQuestionsubmitObj.JudgeInfo = questionSubmitObj.JudgeInfo
	ok, err := myrpc.RpcQuestionSubmitClientImpl.UpdateById(ctx, rpcQuestionsubmitObj)
	if err != nil || !ok.Result {
		myresq.Abort(beeCtx, myresq.OPERATION_ERROR, "题目提交表:更新状态失败")
	}
}

func UpdateQuestionAcceptedNumAdd1(ctx context.Context, beeCtx *beeContext.Context, rpcQuestionObj *rpc_api.RpcQuestionObj) {
	ok, err := myrpc.RpcQuestionClientImpl.Add1AcceptedNum(ctx, rpcQuestionObj)
	if err != nil || !ok.Result {
		myresq.Abort(beeCtx, myresq.OPERATION_ERROR, "题目表:更新题目通过数失败")
	}
}

func DoJudge(beeCtx *beeContext.Context, questionsubmitId int64) *entity.QuestionSubmit {
	var questionSubmitObj entity.QuestionSubmit
	ctx := context.Background()
	// 1）传入题目的提交 id，获取到对应的题目、提交信息（包含代码、编程语言等）
	rpcQuestionsubmitObj, err := myrpc.RpcQuestionSubmitClientImpl.GetById(ctx, &rpc_api.QuestionSubmitGetByIdReq{QuestionSubmitId: questionsubmitId})
	if err != nil {
		myresq.Abort(beeCtx, myresq.NOT_FOUND_ERROR, "提交信息不存在")
		return &questionSubmitObj
	}

	utils.CopyStructFields(*rpcQuestionsubmitObj, &questionSubmitObj)

	// 判断编程语言是否被支持
	if !utils.CheckSame[string]("判题前检查编程语言是否为go", rpcQuestionsubmitObj.Language, questionsubmitlanguageenum.GOLANG.GetValue()) {
		questionSubmitObj.Status = questionsubmitstatusenum.FAILED.GetValue()
		judgeInfoStr, _ := utils.JsonMarshal(model.JudgeInfo{
			Message: judgeinfomessageenum.LANGUAGE_UNSUPPORTED.GetValue(),
			Detail:  "该编程语言暂不支持被判题",
		}, judgeinfomessageenum.LANGUAGE_UNSUPPORTED.GetText())
		questionSubmitObj.JudgeInfo = judgeInfoStr
		UpdateQuestionSubmitObj(ctx, beeCtx, questionSubmitObj, rpcQuestionsubmitObj)
		myresq.Abort(beeCtx, myresq.UNSUPPORTED_ERROR, "该编程语言暂不支持被判题")
		return &questionSubmitObj
	}

	rpcQuestionObj, err := myrpc.RpcQuestionClientImpl.GetById(ctx, &rpc_api.QuestionGetByIdReq{QuestionId: rpcQuestionsubmitObj.QuestionId})
	if err != nil {
		myresq.Abort(beeCtx, myresq.NOT_FOUND_ERROR, "题目不存在")
		return &questionSubmitObj
	}
	var questionObj entity.Question
	utils.CopyStructFields(*rpcQuestionObj, &questionObj)

	// 2）如果题目在判题系统中的处理状态不是”等待中“，就不用重复执行了（在后端提交判题前会修改为为”等待中“）
	if !utils.CheckSame[int32]("判断题目提交状态是否为等待中", questionSubmitObj.Status, questionsubmitstatusenum.WAITING.GetValue()) {
		myresq.Abort(beeCtx, myresq.OPERATION_ERROR, "题目正在判题中")
		return &questionSubmitObj
	}

	// 3）修改题目在判题系统中的处理状态为”判题中“，防止重复执行
	questionSubmitObj.Status = questionsubmitstatusenum.RUNNING.GetValue()
	UpdateQuestionSubmitObj(ctx, beeCtx, questionSubmitObj, rpcQuestionsubmitObj)

	// 获取输入用例
	JudgeCaseStr := questionObj.JudgeCase
	var judgeCaseList []question.JudgeCase
	if err := json.Unmarshal([]byte(JudgeCaseStr), &judgeCaseList); err != nil {
		questionSubmitObj.Status = questionsubmitstatusenum.FAILED.GetValue()
		judgeInfoStr, _ := utils.JsonMarshal(model.JudgeInfo{
			Message: judgeinfomessageenum.SYSTEM_ERROR.GetValue(),
			Detail:  err.Error(),
		}, judgeinfomessageenum.SYSTEM_ERROR.GetText())
		questionSubmitObj.JudgeInfo = judgeInfoStr
		UpdateQuestionSubmitObj(ctx, beeCtx, questionSubmitObj, rpcQuestionsubmitObj)
		myresq.Abort(beeCtx, myresq.OPERATION_ERROR, "输入用例转换失败")
		return &questionSubmitObj
	}
	inputList := make([]string, len(judgeCaseList))
	for i, v := range judgeCaseList {
		inputList[i] = v.Input
	}

	// 4）调用沙箱，获取到执行结果
	codesandboxImpl := codesandbox.CodeSandboxFactory("remote")
	executeCodeResponse, err := codesandbox.CodeSandboxProxy{CodeSandbox: codesandboxImpl}.ExecuteCode(model.ExecuteCodeRequest{
		InputList: inputList,
		Code:      questionSubmitObj.Code,
		Language:  questionSubmitObj.Language,
	})
	if err != nil {
		questionSubmitObj.Status = questionsubmitstatusenum.FAILED.GetValue()

		var e myerror.ErrRemoteSandbox
		if errors.As(err, &e) {
			mylog.Log.Error("代码沙箱返回错误, err=", err.Error())
			judgeInfoStr, _ := utils.JsonMarshal(model.JudgeInfo{
				Message: judgeinfomessageenum.SANDBOX_SYSTEM_ERROR.GetValue(),
				Detail:  err.Error(),
				Memory:  executeCodeResponse.JudgeInfo.Memory,
				Time:    executeCodeResponse.JudgeInfo.Time,
			}, judgeinfomessageenum.SANDBOX_SYSTEM_ERROR.GetText())
			questionSubmitObj.JudgeInfo = judgeInfoStr
			UpdateQuestionSubmitObj(ctx, beeCtx, questionSubmitObj, rpcQuestionsubmitObj)
			myresq.Abort(beeCtx, myresq.OPERATION_ERROR, err.Error())
			return &questionSubmitObj
		}

		judgeInfoStr, _ := utils.JsonMarshal(model.JudgeInfo{
			Message: judgeinfomessageenum.SYSTEM_ERROR.GetValue(),
			Detail:  err.Error(),
			Memory:  executeCodeResponse.JudgeInfo.Memory,
			Time:    executeCodeResponse.JudgeInfo.Time,
		}, judgeinfomessageenum.SYSTEM_ERROR.GetText())
		questionSubmitObj.JudgeInfo = judgeInfoStr
		UpdateQuestionSubmitObj(ctx, beeCtx, questionSubmitObj, rpcQuestionsubmitObj)
		myresq.Abort(beeCtx, myresq.OPERATION_ERROR, err.Error())
		return &questionSubmitObj
	}

	// 5）根据沙箱的执行结果，设置题目的判题状态和信息
	judgeContext := strategy.JudgeContext{
		ExecuteCodeResponse: executeCodeResponse,
		InputList:           inputList,
		JudgeCaseList:       judgeCaseList, //判题用例
		Question:            questionObj,
		QuestionSubmit:      questionSubmitObj, //为了拿到对应的语言，作为判题策略选择依据
	}
	// 使用对应的判题策略进行判题
	judgeInfoResponse := JudgeManager{}.DoJudge(judgeContext)

	// 如果判题结果通过，则修改题目的通过数
	if utils.CheckSame[string]("判断判题结果是否为通过", judgeInfoResponse.Message, judgeinfomessageenum.ACCEPTED.GetValue()) {
		UpdateQuestionAcceptedNumAdd1(ctx, beeCtx, rpcQuestionObj)
	}

	// 更新题目的判题结果judgeInfo到数据库中
	judgeInfoResponseStr, err := utils.JsonMarshal(judgeInfoResponse, "")
	if err != nil {
		questionSubmitObj.Status = questionsubmitstatusenum.FAILED.GetValue()
		judgeInfoStr, _ := utils.JsonMarshal(model.JudgeInfo{
			Message: judgeinfomessageenum.SYSTEM_ERROR.GetValue(),
			Detail:  err.Error(),
		}, judgeinfomessageenum.SYSTEM_ERROR.GetText())
		questionSubmitObj.JudgeInfo = judgeInfoStr
		UpdateQuestionSubmitObj(ctx, beeCtx, questionSubmitObj, rpcQuestionsubmitObj)
		myresq.Abort(beeCtx, myresq.OPERATION_ERROR, "判题信息转换失败")
		return &questionSubmitObj
	}
	questionSubmitObj.JudgeInfo = string(judgeInfoResponseStr)

	// 6）修改题目在判题系统中的处理状态为”成功“
	questionSubmitObj.Status = questionsubmitstatusenum.SUCCEED.GetValue()
	UpdateQuestionSubmitObj(ctx, beeCtx, questionSubmitObj, rpcQuestionsubmitObj)

	return &questionSubmitObj

}
