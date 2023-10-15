<!--
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-14 20:55:53
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-15 18:46:19
 * @FilePath: /xoj-judge-service/README.md
-->
# xoj-judge-service（在线判题系统-判题服务）

## 项目的核心业务

负责将用户提交的代码交给代码沙箱执行后，判断执行结果是否正确、是否满足题目要求，然后将判题结果通过RPC远程调用发送给后端服务。

## 项目本地启动

⚠️ 注意：项目内使用了rpc远程调用，依赖 注册中心已启动、接口提供方已启动(具体见下面《关于 RPC 远程调用》的说明。

* 前置需要：该项目使用了Redis，需要先确保启动了Redis服务

1. 修改/conf 下的配置
    * appconfig.json：修改 `redis` 的连接地址
    * appdynamicconfig.json：修改 `remoteCodeSandboxHost` 为部署代码沙箱服务的服务器IP地址
    * dubbogo.yaml：修改 `nacos` 的连接地址
2. 启动项目
    ```cmd
    go mod tidy
    go run main.go
    ```

## 运行项目中的单元测试

```bash
go test -v ./test
go clean -testcache //清除测试缓存
```

## 关于 RPC 远程调用

该项目内的部分业务使用了dubbo-go 框架的rpc远程调用模式。

* 该项目角色是调用方（Consumer），依赖的提供方（Provide）是[xoj-backend 项目](https://github.com/xiaoxiongmao5/xoj-backend)

* 配置文件位置：/conf/dubbogo.yaml

* 具体业务为为：
     * 获得题目信息 `Question.GetById`
     * 更新题目通过数+1 `Question.Add1AcceptedNum`
     * 获取提交题目信息 `QuestionSubmit.GetById`
     * 更新提交题目信息 `QuestionSubmit.UpdateById`
