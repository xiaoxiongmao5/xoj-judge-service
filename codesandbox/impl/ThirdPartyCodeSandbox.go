/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 12:29:56
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 17:59:33
 */
package impl

import (
	"errors"

	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"
)

type ThirdPartyCodeSandbox struct {
}

func (this ThirdPartyCodeSandbox) ExecuteCode(executeCodeRequest model.ExecuteCodeRequest) (model.ExecuteCodeResponse, error) {
	return model.ExecuteCodeResponse{}, errors.New("暂不支持第三方代码沙箱")
}
