package cron

import (
	"fmt"
	"io"
	"net/http"

	"github.com/heigelove/cpay-payment/internal/repository/mysql/cron_task"
	"go.uber.org/zap"

	"github.com/jakecoffman/cron"
)

func (s *server) AddJob(task *cron_task.CronTask) cron.FuncJob {
	return func() {
		s.taskCount.Add()
		defer s.taskCount.Done()

		// 将 task 信息写入到 Kafka Topic 中，任务执行器订阅 Topic 如果为符合条件的任务并进行执行，反之不执行
		// 为了便于演示，不写入到 Kafka 中，仅记录日志

		msg := fmt.Sprintf("执行任务：(%d)%s [%s]", task.Id, task.Name, task.Spec)
		s.logger.Info(msg)

		if task.Protocol == 2 {
			// 处理 Protocol 2 的任务
			// 执行http请求
			httpClient := &http.Client{}
			req, err := http.NewRequest("GET", task.Command, nil)
			if err != nil {
				s.logger.Error("创建请求失败", zap.Error(err))
				return
			}

			resp, err := httpClient.Do(req)
			if err != nil {
				s.logger.Error("执行请求失败", zap.Error(err))
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				s.logger.Error("请求失败", zap.Int("status", resp.StatusCode))
				return
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				s.logger.Error("读取响应失败", zap.Error(err))
				return
			}
			msg = fmt.Sprintf("任务执行成功：(%d)%s [%s] [%s]", task.Id, task.Name, task.Spec, body)
			s.logger.Info(msg)
		}

	}
}
