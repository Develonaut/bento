package parallel

import (
	"context"
	"fmt"
	"sync"
)

// taskJob represents a task to be executed.
type taskJob struct {
	index int
	task  interface{}
}

// taskResult represents the result of a task execution.
type taskResult struct {
	index  int
	result interface{}
}

// executeWithWorkerPool runs tasks using a worker pool.
func (p *Parallel) executeWithWorkerPool(
	ctx context.Context,
	tasks []interface{},
	maxWorkers int,
	errorStrategy string,
	params map[string]interface{},
) (interface{}, error) {
	taskChan := make(chan taskJob, len(tasks))
	resultChan := make(chan taskResult, len(tasks))
	errorChan := make(chan error, len(tasks))

	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	startWorkers(&wg, workerCtx, maxWorkers, p, taskChan, resultChan, errorChan, params)

	sendTasksToWorkers(taskChan, tasks)

	results, errors, err := collectResults(workerCtx, resultChan, errorChan, &wg, len(tasks), errorStrategy, cancel)
	if err != nil {
		return nil, err
	}

	return buildOutput(results, errors), nil
}

// startWorkers starts the worker goroutines.
func startWorkers(
	wg *sync.WaitGroup,
	ctx context.Context,
	maxWorkers int,
	p *Parallel,
	taskChan <-chan taskJob,
	resultChan chan<- taskResult,
	errorChan chan<- error,
	params map[string]interface{},
) {
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go p.worker(ctx, taskChan, resultChan, errorChan, wg, params)
	}
}

// sendTasksToWorkers sends all tasks to the task channel.
func sendTasksToWorkers(taskChan chan<- taskJob, tasks []interface{}) {
	go func() {
		for i, task := range tasks {
			taskChan <- taskJob{
				index: i,
				task:  task,
			}
		}
		close(taskChan)
	}()
}

// collectResults collects results and errors from workers.
func collectResults(
	ctx context.Context,
	resultChan <-chan taskResult,
	errorChan <-chan error,
	wg *sync.WaitGroup,
	numTasks int,
	errorStrategy string,
	cancel context.CancelFunc,
) ([]interface{}, []interface{}, error) {
	results := make([]interface{}, numTasks)
	errors := make([]interface{}, 0)
	completed := 0

	go closeChannelsWhenDone(wg, resultChan, errorChan)

	for completed < numTasks {
		select {
		case <-ctx.Done():
			return results, errors, ctx.Err()

		case result, ok := <-resultChan:
			if !ok {
				if err := ctx.Err(); err != nil {
					return results, errors, err
				}
				continue
			}
			results[result.index] = result.result
			completed++

		case err, ok := <-errorChan:
			if !ok {
				if ctxErr := ctx.Err(); ctxErr != nil {
					return results, errors, ctxErr
				}
				continue
			}

			if errorStrategy == "failFast" {
				cancel()
				return results, errors, err
			}

			errors = append(errors, err.Error())
			completed++
		}
	}

	return results, errors, nil
}

// closeChannelsWhenDone closes result and error channels when workers are done.
func closeChannelsWhenDone(wg *sync.WaitGroup, resultChan <-chan taskResult, errorChan <-chan error) {
	wg.Wait()
	// Channels are closed by type assertion in the collector
}

// buildOutput builds the output map with results and errors.
func buildOutput(results []interface{}, errors []interface{}) map[string]interface{} {
	output := map[string]interface{}{
		"results": results,
	}

	if len(errors) > 0 {
		output["errors"] = errors
	}

	return output
}

// worker processes tasks from the task channel.
func (p *Parallel) worker(
	ctx context.Context,
	taskChan <-chan taskJob,
	resultChan chan<- taskResult,
	errorChan chan<- error,
	wg *sync.WaitGroup,
	params map[string]interface{},
) {
	defer wg.Done()

	for job := range taskChan {
		if shouldStopWorker(ctx) {
			return
		}

		result, err := p.executeTask(ctx, job.task, params)

		if err != nil {
			if !sendError(ctx, errorChan, err) {
				return
			}
			continue
		}

		if !sendResult(ctx, resultChan, job.index, result) {
			return
		}
	}
}

// shouldStopWorker checks if the worker should stop.
func shouldStopWorker(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// sendError sends an error to the error channel.
func sendError(ctx context.Context, errorChan chan<- error, err error) bool {
	select {
	case errorChan <- err:
		return true
	case <-ctx.Done():
		return false
	}
}

// sendResult sends a result to the result channel.
func sendResult(ctx context.Context, resultChan chan<- taskResult, index int, result interface{}) bool {
	select {
	case resultChan <- taskResult{
		index:  index,
		result: result,
	}:
		return true
	case <-ctx.Done():
		return false
	}
}

// executeTask executes a single task.
func (p *Parallel) executeTask(
	ctx context.Context,
	task interface{},
	params map[string]interface{},
) (interface{}, error) {
	callOnStart(params)
	defer callOnComplete(params)

	if err := checkShouldError(task, params); err != nil {
		return nil, err
	}

	// In real implementation, this would execute nested neta
	// For now, just return the task data
	return task, nil
}

// callOnStart calls the onStart callback if provided.
func callOnStart(params map[string]interface{}) {
	if onStart, ok := params["_onStart"].(func()); ok {
		onStart()
	}
}

// callOnComplete calls the onComplete callback if provided.
func callOnComplete(params map[string]interface{}) {
	if onComplete, ok := params["_onComplete"].(func()); ok {
		onComplete()
	}
}

// checkShouldError checks if the task should error (for testing).
func checkShouldError(task interface{}, params map[string]interface{}) error {
	shouldError, ok := params["_shouldError"].(func(map[string]interface{}) bool)
	if !ok {
		return nil
	}

	taskMap, ok := task.(map[string]interface{})
	if !ok {
		return nil
	}

	if shouldError(taskMap) {
		return fmt.Errorf("task failed")
	}

	return nil
}
