package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
//
//nolint:gosec // disable G115
func Run(tasks []Task, n, m int) error {
	// Проверка, если количество воркеров n <= 0
	if n <= 0 {
		return ErrErrorsLimitExceeded
	}

	// Объявлятся переменные
	var errCount int32
	var ignoreError bool
	maxErrCount := int32(m)
	tasksChan := make(chan Task)
	wg := sync.WaitGroup{}

	if maxErrCount <= 0 {
		ignoreError = true
	}

	// Запускаем n воркеров
	for i := 1; i <= n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for task := range tasksChan {
				// Проверяем, не превышен ли лимит ошибок
				if !ignoreError && atomic.LoadInt32(&errCount) >= maxErrCount {
					return
				}
				// Выполняем задачу
				if err := task(); err != nil {
					atomic.AddInt32(&errCount, 1)
				}
			}
		}()
	}

	// Отправляем задачи в канал
	for _, task := range tasks {
		if !ignoreError && atomic.LoadInt32(&errCount) >= maxErrCount {
			break
		}
		tasksChan <- task
	}

	// Закрываем канал и ждем завершения всех горутин
	close(tasksChan)
	wg.Wait()

	// Проверяем, был ли превышен лимит ошибок
	if !ignoreError && errCount >= maxErrCount {
		return ErrErrorsLimitExceeded
	}

	return nil
}
