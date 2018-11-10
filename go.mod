module task_dispatcher

require (
	github.com/go-redis/redis v6.14.2+incompatible
	github.com/go-sql-driver/mysql v1.4.0
	github.com/gohouse/gorose v1.0.4
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9
	github.com/kataras/iris v10.7.0+incompatible
	github.com/kylelemons/go-gypsy v0.0.0-20160905020020-08cad365cd28
	github.com/shirou/gopsutil v2.18.10+incompatible
	github.com/shotdog/quartz v0.0.0-20160524035313-f6b2ef884f97
	github.com/sirupsen/logrus v1.2.0
	github.com/sunpengwei1992/task_dispatcher v0.0.0-20181109145256-fae71bbe29dc
)

replace github.com/sunpengwei1992/task_dispatcher v0.0.0-20181109145256-fae71bbe29dc => /
