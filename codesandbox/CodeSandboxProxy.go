/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-08 16:22:50
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 17:58:22
 */
package codesandbox

import (
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/mylog"
)

type CodeSandboxProxy struct {
	CodeSandbox CodeSandbox
}

func (this CodeSandboxProxy) ExecuteCode(executeCodeRequest model.ExecuteCodeRequest) (model.ExecuteCodeResponse, error) {
	mylog.Log.Infof("代码沙箱请求信息：%v", executeCodeRequest)
	executeCodeResponse, err := this.CodeSandbox.ExecuteCode(executeCodeRequest)
	mylog.Log.Infof("代码沙箱请响应信息：%v", executeCodeResponse)
	return executeCodeResponse, err
}
