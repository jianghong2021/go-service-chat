package cmd

import (
	"fmt"
	"goflylivechat/models"
	"time"
)

var (
	STORE_HOUR int64         = 6
	CLEAR_HOUR time.Duration = 1 * time.Hour
)

func StartCronJobs() {
	ticker := time.NewTicker(CLEAR_HOUR)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				doTask()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func doTask() {
	fmt.Println("消息清除任务执行:", time.Now().Format("2006-01-02 15:04:05"))
	now := time.Now().Unix() - STORE_HOUR*3600
	count := models.DeleteMessage("created_at <= FROM_UNIXTIME(?)", now)
	fmt.Println("清除消息数量:", count, now)
}
