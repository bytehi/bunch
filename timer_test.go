package bunch

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	// 创建一个 Timer 实例
	timer := &Timer{}

	// 添加一个定时任务，1秒后执行
	timer.Add(time.Second, func() {
		fmt.Println(time.Now(), "Task 1 executed")
	})

	// 添加另一个定时任务，2秒后执行
	timer.Add(2*time.Second, func() {
		fmt.Println(time.Now(), "Task 2 executed")
	})

	dead := time.Now().Add(10 * time.Second)
	fmt.Println(dead)
	for now := time.Now(); now.Before(dead); now = time.Now() {
		time.Sleep(time.Millisecond)
		timer.Timeout(time.Now())
	}
	fmt.Println("done")
}
