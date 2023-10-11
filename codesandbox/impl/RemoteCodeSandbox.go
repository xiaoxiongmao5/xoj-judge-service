/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 12:29:56
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-11 16:48:16
 */
package impl

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/config"
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

func (this RemoteCodeSandbox) ExecuteCode(executeCodeRequest model.ExecuteCodeRequest) (executeCodeResponse model.ExecuteCodeResponse, err error) {
	// 将请求数据结构体编码为 JSON 字符串
	requestBody, err := json.Marshal(executeCodeRequest)
	if err != nil {
		msg := fmt.Sprintf("编码请求数据结构体失败：%s", err.Error())
		mylog.Log.Error(msg)
		return executeCodeResponse, errors.New(msg)
	}

	targetURL := config.AppConfigDynamic.RemoteCodeSandboxHost
	bodyBytes, err := utils.SendHTTPRequest(
		"POST",
		targetURL,
		requestBody,
	)
	if err != nil {
		msg := fmt.Sprintf("请求代码沙箱失败：%s", err.Error())
		mylog.Log.Error(msg)
		return executeCodeResponse, errors.New(msg)
	}

	// 解析 JSON
	var responseData ResponseData
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		msg := fmt.Sprintf("解析JSON响应数据失败：%s", err.Error())
		mylog.Log.Error(msg)
		return executeCodeResponse, errors.New(msg)
	}

	if responseData.Code != 0 {
		msg := fmt.Sprintf("响应数据Code校验失败：%s", err.Error())
		mylog.Log.Error(msg)
		return executeCodeResponse, errors.New(msg)
	}

	utils.CopyStructFields(responseData.Data, &executeCodeResponse)

	return executeCodeResponse, nil
}
