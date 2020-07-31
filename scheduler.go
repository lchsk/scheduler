package scheduler

import (
	"time"
)

type Task struct {
	Every         time.Duration
	LastExecution time.Time

	Func func()
}

type TaskManager struct {
	Tasks    []*Task
	Interval time.Duration
}

func (mgr *TaskManager) Check() {
	for _, task := range mgr.Tasks {
		current := time.Now().UTC()
		diff := current.Sub(task.LastExecution)

		if diff > task.Every {
			if task.Func != nil {
				go task.Func()
			}
			task.LastExecution = current
		}
	}
}

func (mgr *TaskManager) Wait() {
	wait := make(chan bool)

	go func() {
		for {
			mgr.Check()
			time.Sleep(mgr.Interval)
		}
	}()

	<-wait
}

func (mgr *TaskManager) Schedule(every time.Duration, fun func()) {
	mgr.Tasks = append(mgr.Tasks, &Task{
		Every:         every,
		LastExecution: time.Now().UTC(),
		Func:          fun,
	})

	mgr.Interval = 999 * time.Hour

	for _, task := range mgr.Tasks {
		if task.Every < mgr.Interval {
			mgr.Interval = task.Every
		}
	}
}
