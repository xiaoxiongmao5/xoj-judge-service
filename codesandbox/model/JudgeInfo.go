/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-09-29 09:12:38
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-15 17:48:27
 */
package model

type JudgeInfo struct {
	// 程序执行信息(这是由判题系统写的)
	Message string `json:"message"`
	// 程序执行信息详细(这是由判题系统写的)
	Detail string `json:"detail"`
	// 消耗内存
	Memory int64 `json:"memory"`
	// 消耗时间
	Time int64 `json:"time"`
}
