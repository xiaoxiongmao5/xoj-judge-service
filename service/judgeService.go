/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 13:27:42
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-10 15:48:32
 */
package service

import (
	"context"
	"encoding/json"

	beeContext "github.com/beego/beego/v2/server/web/context"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/dto/question"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/entity"
	judgeinfomessageenum "github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/enums/JudgeInfoMessageEnum"
	questionsubmitstatusenum "github.com/xiaoxiongmao5/xoj/xoj-judge-service/model/enums/QuestionSubmitStatusEnum"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/myresq"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/myrpc"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/rpc_api"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/service/strategy"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/utils"
)

func UpdateQuestionSubmitObjStatus(beeCtx *beeContext.Context, rpcQuestionsubmitObj *rpc_api.RpcQuestionSubmitObj) {
	ok, err := myrpc.RpcQuestionSubmitClientImpl.UpdateById(context.Background(), rpcQuestionsubmitObj)
	if err != nil || !ok.Result {
		myresq.Abort(beeCtx, myresq.OPERATION_ERROR, "题目提交表:更新状态失败")
	}
}

func UpdateQuestionAcceptedNumAdd1(beeCtx *beeContext.Context, rpcQuestionObj *rpc_api.RpcQuestionObj) {
	ok, err := myrpc.RpcQuestionClientImpl.Add1AcceptedNum(context.Background(), rpcQuestionObj)
	if err != nil || !ok.Result {
		myresq.Abort(beeCtx, myresq.OPERATION_ERROR, "题目表:更新题目通过数失败")
	}
}

func DoJudge(beeCtx *beeContext.Context, questionsubmitId int64) *entity.QuestionSubmit {
	// 1）传入题目的提交 id，获取到对应的题目、提交信息（包含代码、编程语言等）
	rpcQuestionsubmitObj, err := myrpc.RpcQuestionSubmitClientImpl.GetById(context.Background(), &rpc_api.QuestionSubmitGetByIdReq{QuestionSubmitId: questionsubmitId})
	if err != nil {
		myresq.Abort(beeCtx, myresq.NOT_FOUND_ERROR, "提交信息不存在")
		return nil
	}
	var questionSubmitObj entity.QuestionSubmit
	utils.CopyStructFields(*rpcQuestionsubmitObj, &questionSubmitObj)

	rpcQuestionObj, err := myrpc.RpcQuestionClientImpl.GetById(context.Background(), &rpc_api.QuestionGetByIdReq{QuestionId: rpcQuestionsubmitObj.QuestionId})
	if err != nil {
		myresq.Abort(beeCtx, myresq.NOT_FOUND_ERROR, "题目不存在")
		return nil
	}
	var questionObj entity.Question
	utils.CopyStructFields(*rpcQuestionObj, &questionObj)

	// 2）如果题目在判题系统中的处理状态不是”等待中“，就不用重复执行了（在后端提交判题前会修改为为”等待中“）
	if !utils.CheckSame[int32]("判断题目提交状态是否为等待中", questionSubmitObj.Status, questionsubmitstatusenum.WAITING.GetValue()) {
		myresq.Abort(beeCtx, myresq.OPERATION_ERROR, "题目正在判题中")
		return nil
	}

	// 3）修改题目在判题系统中的处理状态为”判题中“，防止重复执行
	questionSubmitObj.Status = questionsubmitstatusenum.RUNNING.GetValue()
	rpcQuestionsubmitObj.Status = questionSubmitObj.Status
	UpdateQuestionSubmitObjStatus(beeCtx, rpcQuestionsubmitObj)

	// 获取输入用例
	JudgeCaseStr := questionObj.JudgeCase
	var judgeCaseList []question.JudgeCase
	if err := json.Unmarshal([]byte(JudgeCaseStr), &judgeCaseList); err != nil {
		questionSubmitObj.Status = questionsubmitstatusenum.FAILED.GetValue()
		rpcQuestionsubmitObj.Status = questionSubmitObj.Status
		UpdateQuestionSubmitObjStatus(beeCtx, rpcQuestionsubmitObj)
		myresq.Abort(beeCtx, myresq.OPERATION_ERROR, "输入用例转换失败")
		return nil
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
		rpcQuestionsubmitObj.Status = questionSubmitObj.Status
		UpdateQuestionSubmitObjStatus(beeCtx, rpcQuestionsubmitObj)
		myresq.Abort(beeCtx, myresq.OPERATION_ERROR, err.Error())
		return nil
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
		UpdateQuestionAcceptedNumAdd1(beeCtx, rpcQuestionObj)
	}

	// 更新题目的判题结果judgeInfo到数据库中
	if judgeInfoResponseStr, err := json.Marshal(judgeInfoResponse); err != nil {
		questionSubmitObj.Status = questionsubmitstatusenum.FAILED.GetValue()
		rpcQuestionsubmitObj.Status = questionSubmitObj.Status
		UpdateQuestionSubmitObjStatus(beeCtx, rpcQuestionsubmitObj)
		myresq.Abort(beeCtx, myresq.OPERATION_ERROR, "判题信息转换失败")
		return nil
	} else {
		questionSubmitObj.JudgeInfo = string(judgeInfoResponseStr)
		rpcQuestionsubmitObj.JudgeInfo = questionSubmitObj.JudgeInfo
	}

	// 6）修改题目在判题系统中的处理状态为”成功“
	questionSubmitObj.Status = questionsubmitstatusenum.SUCCEED.GetValue()
	rpcQuestionsubmitObj.Status = questionSubmitObj.Status
	UpdateQuestionSubmitObjStatus(beeCtx, rpcQuestionsubmitObj)

	return &questionSubmitObj

}
