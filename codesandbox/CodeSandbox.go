/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 12:14:47
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 18:00:39
 */
package codesandbox

import "github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/model"

// 代码沙箱
type CodeSandbox interface {
	ExecuteCode(executeCodeRequest model.ExecuteCodeRequest) (model.ExecuteCodeResponse, error)
}
