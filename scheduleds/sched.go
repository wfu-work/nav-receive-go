package scheduleds

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func Init() {
	timers := NewTimerTask()
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())
		_, _ = timers.AddTaskByFunc("Rtlogging", "0 */5 * * * *", func() {
			err := RtloggingSched()
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "每五分钟定时检查设备状态", option...)
	}()
}
