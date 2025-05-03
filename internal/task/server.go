package task

import (
	"awesomeProject/pkg/configs"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

func InitTaskServer() {

	redisAddr := configs.GetConfig().Redis.Host + ":" + configs.GetConfig().Redis.Port

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			// critical occupy 60% of your total processing power.
			// default occupy 30%.
			// low occupy 10%.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	var mux = asynq.NewServeMux()

	newContainerProcessor()
	processor.Register(mux)

	if err := srv.Run(mux); err != nil {
		logrus.Fatal(err)
	}
}
