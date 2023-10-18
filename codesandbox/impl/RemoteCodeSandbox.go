/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 12:29:56
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-18 00:46:34
 */
package impl

import (
	"encoding/json"

	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/config"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/myerror"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/utils"
)

type RemoteCodeSandbox struct {
}

// 响应的数据结构
type ResponseData struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Data    model.ExecuteCodeResponse `json:"data"`
}

func (this RemoteCodeSandbox) ExecuteCode(executeCodeRequest model.ExecuteCodeRequest) (model.ExecuteCodeResponse, error) {
	tag := "RemoteCodeSandbox ExecuteCode:"
	var executeCodeResponse model.ExecuteCodeResponse
	// 将请求数据结构体编码为 JSON 字符串
	requestBody, err := json.Marshal(executeCodeRequest)
	if err != nil {
		mylog.Log.Errorf("%s json.Marshal(executeCodeRequest) 失败, err=[%s]", tag, err.Error())
		return executeCodeResponse, err
	}

	targetURL := config.AppConfigDynamic.RemoteCodeSandboxHost
	bodyBytes, err := utils.SendHTTPRequest(
		"POST",
		targetURL,
		requestBody,
	)
	if err != nil {
		mylog.Log.Errorf("%s utils.SendHTTPRequest 失败, err=[%s]", tag, err.Error())
		return executeCodeResponse, err
	}

	// 解析 JSON
	var responseData ResponseData
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		mylog.Log.Errorf("%s json.Unmarshal(bodyBytes, &responseData) 失败, err=[%s]", tag, err.Error())
		return executeCodeResponse, myerror.ErrRemoteSandbox{Message: err.Error()}
	}

	utils.CopyStructFields(responseData.Data, &executeCodeResponse)

	if responseData.Code != 0 {
		mylog.Log.Errorf("%s responseData.Code != 0, Code=[%d],Message=[%s]", tag, responseData.Code, responseData.Message)
		return executeCodeResponse, myerror.ErrRemoteSandbox{Message: responseData.Message}
	}

	return executeCodeResponse, nil
}
