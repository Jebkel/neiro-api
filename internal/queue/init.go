package queue

import (
	mailQueue "neiro-api/internal/queue/mail"
	"neiro-api/pkg/queue"
)

func Init() {
	q := queue.GetQueue()

	q.RegisterHandler("mail", mailQueue.TaskHandlerWrapper)
	q.Start()
}
