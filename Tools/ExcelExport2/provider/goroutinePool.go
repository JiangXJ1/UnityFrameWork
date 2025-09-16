package provider

import (
	"sync"
)

/*
 	// 创建一个 Goroutine 池
	pool := NewGoroutinePool(3, 10)

	// 启动 Goroutine 池
	pool.Start()

	// 提交任务
	for i := 1; i <= 10; i++ {
		pool.Submit(&Task{ID: i, Message: fmt.Sprintf("Task %d", i)})
	}

	// 停止 Goroutine 池
	pool.WaitComplete()
*/

// Goroutine 池
type GoroutinePool struct {
	TaskQueue chan *Task
	WorkerNum int
	wg        sync.WaitGroup
}

// 创建 Goroutine 池
// workerNum: 工作者数量
// queueSize: 通道初始容量（任务队列大小）
func NewGoroutinePool(workerNum int, queueSize int) *GoroutinePool {
	return &GoroutinePool{
		TaskQueue: make(chan *Task, queueSize),
		WorkerNum: workerNum,
	}
}

// 启动 Worker
func (p *GoroutinePool) Start() {
	for i := 0; i < p.WorkerNum; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			for task := range p.TaskQueue {
				task.Execute()
			}
		}(i)
	}
}

// WaitComplete 停止加入任务队列并等待所有任务结束
func (p *GoroutinePool) WaitComplete() {
	close(p.TaskQueue) // 关闭任务队列(关闭后无法再提交任务，但已提交的任务会继续执行,全部完成后才会退出)
	p.wg.Wait()        // 等待所有 worker 结束
	//fmt.Println("\nGoroutine pool stopped!")
}

// 提交任务
func (p *GoroutinePool) Submit(task *Task) {
	//fmt.Printf("++++++++++Submitting task: %s %v\n", task.FileName, time.Now())
	p.TaskQueue <- task
}
