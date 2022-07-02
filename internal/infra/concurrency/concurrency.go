package concurrency

import "sync"

type TaskResult struct {
	Result interface{}
	Err    error
	index  int
}

type Task = func() (interface{}, error)
type TaskResultMap map[int]TaskResult

func ExecuteConcurrentTasks(tasks ...Task) TaskResultMap {
	taskResult := make(TaskResultMap, 0)
	asyncTaskChannel := make(chan TaskResult)

	wg := sync.WaitGroup{}
	wg.Add(len(tasks))

	for index, task := range tasks {
		go func(index int, task Task, channel chan TaskResult) {
			result, err := task()
			channel <- TaskResult{
				Result: result,
				Err:    err,
				index:  index,
			}
		}(index, task, asyncTaskChannel)
	}

	for i := 0; i < len(tasks); i++ {
		message := <-asyncTaskChannel
		taskResult[message.index] = message
		wg.Done()
	}

	wg.Wait()

	return taskResult
}
