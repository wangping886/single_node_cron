package single_node_cron

import (
	"github.com/gomodule/redigo/redis"
	"testing"
)

func TestCronTask(t *testing.T) {
	//传入redis 连接池
	cron := NewSingleNodeCron(redis.Pool{}, "project_name")
	cron.AddFunc("* * * * * *", "method_name", ExecMethod)
	cron.Start()
	//signal
	//cron.Stop()
}
func ExecMethod() {
	//your process logic
}
