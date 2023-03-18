package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

func WorkerInIt() {
	srv := asynq.NewServer(asynq.RedisClientOpt{Addr: "127.0.0.1:6379"}, asynq.Config{Concurrency: 4})

	mux := asynq.NewServeMux()
	mux.HandleFunc(string(SendTask), SendTaskMessage)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}

func SendTaskMessage(ctx context.Context, task *asynq.Task) error {
	var ep EmailPayload
	if err := json.Unmarshal(task.Payload(), &ep); err != nil {
		return err
	}
	fmt.Println("sent task to :", ep.UserID)
	return nil
}
