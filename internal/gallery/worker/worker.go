package worker

import (
	"fmt"
	"image-gallery/internal/gallery/entity"
	"image-gallery/internal/gallery/repo"
	"sync"
)

type Worker struct {
	WorkerCount int
	TaskQueue   chan entity.Image
	Result      chan string
	Out         <-chan string
	Wg          sync.WaitGroup
	repo        repo.Repository
}

func NewWorker(workerCount int, taskQueue chan entity.Image, result chan string, Out <-chan string, repo repo.Repository) *Worker {
	return &Worker{
		WorkerCount: workerCount,
		TaskQueue:   taskQueue,
		Result:      result,
		Out:         Out,
		repo:        repo,
	}
}

func (w *Worker) WorkerRun(id int) {
	defer w.Wg.Done()

	for {
		task, ok := <-w.TaskQueue
		if !ok {
			fmt.Printf("Worker %d shutting down\n", id)
			return
		}
		err := w.repo.CreatePhoto(task)
		if err != nil {
			fmt.Printf("error with creating images with id: %d", task.Id)
		}

		fmt.Printf("Worker %d completed task \n", id)
	}
}
