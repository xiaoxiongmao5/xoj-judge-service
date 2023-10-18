/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-17 11:14:59
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-17 11:35:46
 */
package myerror

type ErrRemoteSandbox struct {
	Message string
}

func (this ErrRemoteSandbox) Error() string {
	return this.Message
}
