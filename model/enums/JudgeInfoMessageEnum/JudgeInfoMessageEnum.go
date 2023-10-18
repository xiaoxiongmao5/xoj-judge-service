/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-09-29 21:42:35
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-17 23:09:13
 * @FilePath: /xoj-backend/model/enums/JudgeInfoMessageEnum.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package judgeinfomessageenum

type JudgeInfoMessageEnum string

func (this JudgeInfoMessageEnum) GetValue() string {
	return string(this)
}

func (this JudgeInfoMessageEnum) GetText() string {
	return JudgeInfoMessageEnumName[this]
}

// 提交题目的判题结果

const (
	ACCEPTED                    JudgeInfoMessageEnum = "Accepted"
	WRONG_ANSWER                JudgeInfoMessageEnum = "Wrong Answer"
	COMPILE_ERROR               JudgeInfoMessageEnum = "Compile Error"
	MEMORY_LIMIT_EXCEEDED       JudgeInfoMessageEnum = "Memory Limit Exceeded"
	OUT_OF_MEMORY               JudgeInfoMessageEnum = "Out Of Memory"
	COMPILE_TIME_LIMIT_EXCEEDED JudgeInfoMessageEnum = "Compile Time Limit Exceeded"
	RUN_TIME_LIMIT_EXCEEDED     JudgeInfoMessageEnum = "Run Time Limit Exceeded"
	TIME_LIMIT_EXCEEDED         JudgeInfoMessageEnum = "Time Limit Exceeded"
	PRESENTATION_ERROR          JudgeInfoMessageEnum = "Presentation Error"
	WAITING                     JudgeInfoMessageEnum = "Waiting"
	OUTPUT_LIMIT_EXCEEDED       JudgeInfoMessageEnum = "Output Limit Exceeded"
	DANGEROUS_OPERATION         JudgeInfoMessageEnum = "Dangerous Operation"
	RUNTIME_ERROR               JudgeInfoMessageEnum = "Runtime Error"
	SYSTEM_ERROR                JudgeInfoMessageEnum = "System Error"
	LANGUAGE_UNSUPPORTED        JudgeInfoMessageEnum = "Language UnSupported"
	SANDBOX_SYSTEM_ERROR        JudgeInfoMessageEnum = "Sandbox System Error"
)

var JudgeInfoMessageEnumName = map[JudgeInfoMessageEnum]string{
	ACCEPTED:                    "成功",
	WRONG_ANSWER:                "答案错误",
	COMPILE_ERROR:               "编译错误",
	MEMORY_LIMIT_EXCEEDED:       "内存溢出",
	OUT_OF_MEMORY:               "内存不足",
	COMPILE_TIME_LIMIT_EXCEEDED: "编译超时",
	RUN_TIME_LIMIT_EXCEEDED:     "运行超时",
	TIME_LIMIT_EXCEEDED:         "超题限时",
	PRESENTATION_ERROR:          "展示错误",
	WAITING:                     "等待中",
	OUTPUT_LIMIT_EXCEEDED:       "输出溢出",
	DANGEROUS_OPERATION:         "危险操作",
	RUNTIME_ERROR:               "运行错误",
	SYSTEM_ERROR:                "系统错误",
	LANGUAGE_UNSUPPORTED:        "语言不支持",
	SANDBOX_SYSTEM_ERROR:        "沙箱系统错误",
}

// 根据 value 获取枚举
func GetEnumByValue(value string) {

}
