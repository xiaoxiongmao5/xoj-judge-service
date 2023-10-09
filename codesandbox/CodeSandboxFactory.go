/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-02 12:14:47
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 11:20:59
 */
package codesandbox

import "github.com/xiaoxiongmao5/xoj/xoj-judge-service/codesandbox/impl"

// 代码沙箱工厂（根据字符串参数创建指定的代码沙箱实例）
func CodeSandboxFactory(codesandboxType string) CodeSandbox {
	switch codesandboxType {
	case "example":
		return impl.ExampleCodeSandbox{}
	case "remote":
		return impl.RemoteCodeSandbox{}
	case "thirdParty":
		return impl.ThirdPartyCodeSandbox{}
	default:
		return impl.ExampleCodeSandbox{}
	}
}
