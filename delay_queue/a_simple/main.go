package main

import (
	"fmt"
	"time"
)

/*
基于go实现延迟队列
*/
type Task struct {
	ExecuteTime time.Time
	Job         func()
}

type DelayQueue struct {
	Tasks []*Task
}

func (d *DelayQueue) AddTask(t *Task) {
	d.Tasks = append(d.Tasks, t)
}

func (d *DelayQueue) RemoveTask() {
	//FIFO: remove the first task to enqueue
	d.Tasks = d.Tasks[1:]
}

func (d *DelayQueue) ExecuteTask() {
	for len(d.Tasks) > 0 {
		//dequeue a task
		currentTask := d.Tasks[0]
		if time.Now().Before(currentTask.ExecuteTime) {
			//if the task execution time is not up, wait
			time.Sleep(currentTask.ExecuteTime.Sub(time.Now()))
		}
		//execute the task
		currentTask.Job()
		//remove task who has been executed
		d.RemoveTask()
	}

}

func main() {
	fmt.Println("start delayQueue")
	delayQueue := &DelayQueue{}
	firstTask := &Task{
		ExecuteTime: time.Now().Add(time.Second * 1),
		Job: func() {
			fmt.Println("executed task 1 after delay")
		},
	}
	delayQueue.AddTask(firstTask)
	secondTask := &Task{
		ExecuteTime: time.Now().Add(time.Second * 7),
		Job: func() {
			fmt.Println("executed task 2 after delay")
		},
	}
	delayQueue.AddTask(secondTask)
	delayQueue.ExecuteTask()
	fmt.Println("all tasks have been done!!!")
}
