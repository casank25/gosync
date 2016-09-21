package main

type Dispatcher struct {
	WorkerPool chan chan Job
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool}
}

func (d *Dispatcher) Run() {
	for i := 1; i <= options.Workers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}
	go d.Dispatch()
}

func (d *Dispatcher) Dispatch() {
	for {
		select {
		case job := <-JobQueue:
			//print something here
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		}
	}
}
