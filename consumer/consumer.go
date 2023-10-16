/*
 * @Author: 小熊 627516430@qq.com
 * @Date: 2023-10-10 15:51:10
 * @LastEditors: 小熊 627516430@qq.com
 * @LastEditTime: 2023-10-16 11:23:58
 * @FilePath: /xoj-judge-service/consumer/consumer.go
 * @Description: 消费者
 */
package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/mylog"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/service"
	"github.com/xiaoxiongmao5/xoj/xoj-judge-service/utils"
)

const (
	QuestionSubmit2QueueKey = "Queue_QuestionSubmit"
)

func PopQuestionSubmit2Queue(ctx context.Context, client *redis.Client, mu *sync.Mutex) {

	for {
		select {
		case <-ctx.Done():
			return // 如果上下文被取消，则退出协程
		default:
			if client == nil {
				mylog.Log.Error("redis client is nil")
				continue
			}
			// 使用互斥锁保护对队列的访问
			mu.Lock()
			// 从队列的左侧（头部）获取消息，这将阻塞等待直到有消息可用
			message, err := client.LPop(ctx, QuestionSubmit2QueueKey).Result()
			mu.Unlock()

			if err == redis.Nil {
				// 队列为空，等待一段时间后继续尝试
				time.Sleep(time.Second)
				continue
			} else if err != nil {
				mylog.Log.Error("无法获取消息:", err)
				return
			}
			mylog.Log.Infof("已消费消息: %s", message)

			questionsubmitId, err := utils.String2Int64(message)
			if err != nil {
				mylog.Log.Errorf("%s 转提交题目Id失败", message)
				continue
			}

			// 执行判题
			service.DoJudge(nil, questionsubmitId)
		}
	}
}
