package single_node_cron

import (
	"github.com/garyburd/redigo/redis"
	"gopkg.in/robfig/cron.v2"
	"log"
)

type singleNodeCron struct {
	*cron.Cron
	redisPool *redis.Pool
	prefix    string
}

func NewSingleNodeCron(pool *redis.Pool, prefix string) *singleNodeCron {
	return &singleNodeCron{
		redisPool: pool,
		Cron:      cron.New(),
		prefix:    prefix,
	}
}

func (c *singleNodeCron) AddFunc(spec, subLockName string, cmd func()) (cron.EntryID, error) {
	return c.Cron.AddFunc(spec, func() {
		if getRedisLock(c.redisPool, c.prefix+subLockName) {
			cmd()
		}
	})
}

func (c *singleNodeCron) AddJob(spec, subLockName string, cmd cron.Job) (cron.EntryID, error) {
	return c.Cron.AddJob(spec, cron.FuncJob(func() {
		if getRedisLock(c.redisPool, c.prefix+subLockName) {
			cmd.Run()
		}
	}))
}

func (c *singleNodeCron) Start() {
	c.Cron.Start()
}

func (c *singleNodeCron) Stop() {
	c.Cron.Stop()
}

func getRedisLock(redisPool *redis.Pool, lockKey string) bool {
	conn := redisPool.Get()
	defer conn.Close()
	// SET resource-name anystring NX EX max-lock-time  原子操作 脚本级别的分布式锁
	// 适用分布式环境下需要控制单节点执行任务
	result, err := redis.String(conn.Do("SET", lockKey, "lock", "EX", 3, "NX"))
	switch err {
	case nil:
		if result != "OK" && result != "ok" {
			return false
		}
	case redis.ErrNil:
		return false
	default:
		log.Println("msg", "distributedLockError", "err", err, "key", lockKey)
		return false
	}

	return true
}
