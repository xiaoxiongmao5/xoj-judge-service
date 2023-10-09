/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-09 13:30:41
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-09 13:31:25
 */
package swagtype

type CommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
