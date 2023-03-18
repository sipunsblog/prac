package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type EmailPayload struct {
	UserID string
}

type TaskName string

const (
	SendTask TaskName = "task:sendTask"
	DropMsg  TaskName = "task:dropmsg"
)

func CreateClient() {
	client := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr: "127.0.0.1:6379",
		},
	)
	pl, err := json.Marshal(EmailPayload{UserID: "saumya@gamil.com"})
	if err != nil {
		fmt.Printf("Error marshalling email payload")
	}
	sendTask := asynq.NewTask(string(SendTask), pl)
	dropMsg := asynq.NewTask(string(DropMsg), pl)

	info1, err := client.Enqueue(sendTask, asynq.ProcessIn(10*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success1", info1)

	_, _ = client.Enqueue(dropMsg, asynq.ProcessIn(10*time.Second))

	fmt.Println("Success2", info1)

}
