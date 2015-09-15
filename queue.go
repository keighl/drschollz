package drschollz

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"
)

// Job packages relevant values for an error notification
type Job struct {
	// The raw error value
	Err       error
	Time      time.Time
	BackTrace string
	// A list of abitrary extra values to be included in the email
	Extras []interface{}
}

// Error is the standard error interface
func (j Job) Error() string {
	return j.Err.Error()
}

func (j Job) subject() string {
	if len(j.Error()) <= 60 {
		return j.Error()
	}
	return j.Error()[0:60] + "..."
}

// AppName is a convenience method for accessing the global AppName within email templates
func (j Job) AppName() string {
	return Conf.AppName
}

/////////////////

// Queue - a pool of workers to handle asynchronous error notifications
type Queue struct {
	Running bool
	Log     bool

	workers    []worker
	jobQueue   chan Job
	quitChan   chan *sync.WaitGroup
	workerPool chan worker
	logger     *log.Logger
}

// Start generates a new Queue and starts its worker pool
func Start(workers int) (*Queue, error) {
	q := &Queue{}
	err := q.Start(workers)
	if err != nil {
		return nil, err
	}

	return q, nil
}

// Println will print to the logger when q.Log is true
func (q *Queue) Println(v ...interface{}) {
	if q.Log {
		q.logger.Println(v...)
	}
}

// Start begins a queues worker pool
func (q *Queue) Start(workers int) error {

	if workers <= 0 {
		return errors.New("Gimme at least one worker!")
	}

	if q.Running {
		return errors.New("We already runnin'!!")
	}

	q.workers = []worker{}
	q.jobQueue = make(chan Job, 100)
	q.workerPool = make(chan worker, workers)
	q.quitChan = make(chan *sync.WaitGroup)
	q.logger = log.New(os.Stdout, "schollz: ", log.Ldate|log.Ltime)
	q.Log = true

	for i := 0; i < workers; i++ {
		worker := worker{
			ID:       i,
			WorkChan: make(chan Job),
			QuitChan: make(chan *sync.WaitGroup),
			Queue:    q,
		}
		worker.start()
		q.workers = append(q.workers, worker)
	}

	go func(q *Queue) {
		for {
			select {
			case job := <-q.jobQueue:
				go func() {
					worker := <-q.workerPool
					worker.WorkChan <- job
				}()
			case wg := <-q.quitChan:
				defer wg.Done()
				for _, w := range q.workers {
					go func(w worker) {
						w.QuitChan <- wg
					}(w)
				}
				return
			}
		}
	}(q)

	q.Running = true
	return nil
}

// Stop will shut down the queue's worker pool. It blocks until all workers are stopped
func (q *Queue) Stop() {
	if !q.Running {
		return
	}
	q.Println("Dr. Schollz is closing up shop...")

	var wg sync.WaitGroup
	wg.Add(1 + len(q.workers))
	go func() {
		q.quitChan <- &wg
	}()
	wg.Wait()
	q.Running = false
}

/////////////////

type worker struct {
	ID       int
	Queue    *Queue
	WorkChan chan Job
	QuitChan chan *sync.WaitGroup
}

func (w worker) start() {
	go func() {
		for {
			w.Queue.workerPool <- w
			select {
			case job := <-w.WorkChan:
				w.Queue.Println("Wrkr", w.ID, "->", job.Error)
				err := Deliver(job)
				if err != nil {
					w.Queue.Println("Wrkr", w.ID, "-> errored out!", err)
				} else {
					w.Queue.Println("Wrkr", w.ID, "-> delivered", job.Error)
				}
			case wg := <-w.QuitChan:
				defer wg.Done()
				return
			}
		}
	}()
}
