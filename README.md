# single_node_cron
封装robfig/cron，用于分布式多节点环境，控制单节点运行任务

通过redis分布式锁控制，
通过redis 的SET resource-name anystring NX EX max-lock-time(redis 2.6.12版本以上) 原子操作 。避免lua实现以及保证原子性
