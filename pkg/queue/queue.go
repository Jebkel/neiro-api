package queue

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Task struct {
	QueueName string
	Data      interface{}
}

type Queue struct {
	tasks    chan Task
	handlers map[string]TaskHandler
	wg       sync.WaitGroup
}

var (
	instance *Queue
	once     sync.Once
)

type TaskHandler func(data interface{}) error

func GetQueue() *Queue {
	once.Do(func() {
		instance = &Queue{
			tasks:    make(chan Task, 100),
			handlers: make(map[string]TaskHandler),
		}
	})
	return instance
}

func (q *Queue) AddTask(queueName string, taskData interface{}) {
	q.tasks <- Task{QueueName: queueName, Data: taskData}
}

func (q *Queue) RegisterHandler(queueName string, handler TaskHandler) {
	q.handlers[queueName] = handler
}

func (q *Queue) Start() {
	for task := range q.tasks {
		log.Debugf("Полученая новая таска: %v", task)
		handler, ok := q.handlers[task.QueueName]
		if !ok {
			log.Errorf("Обработчик для очереди '%s' не найден", task.QueueName)
			continue
		}
		q.wg.Add(1)
		go func(handler TaskHandler, data interface{}) {
			defer q.wg.Done()
			var err error
			for attempt := 0; attempt < 3; attempt++ {
				err = handler(data)
				if err == nil {
					return
				}
				log.Errorf("Ошибка при выполении фоновой задачи: %s. Попытка %d, ошибка: %v", err.Error(), attempt+1, err.Error())
				time.Sleep(5 * time.Second)
			}
			log.Error(err)
		}(handler, task.Data)
	}
}
