package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/abdullahnettoor/tqwp"
)

type CustomTask struct {
	tqwp.TaskModel
	Id   uint
	Data any
}

func (t *CustomTask) Process() error {
	num, isInt := t.Data.(int)
	if !isInt {
		return fmt.Errorf("invalid type")
	}
	divisor := rand.Intn(2)
	if divisor == 0 {
		return fmt.Errorf("division by zero")
	}
	t.Data = num / divisor
	time.Sleep(time.Millisecond * 10)
	return nil
}

func main() {

	wp := tqwp.New(&tqwp.WorkerPoolConfig{MaxRetries: 3, NumOfWorkers: 10})
	defer wp.Summary()
	defer wp.Stop()

	wp.EnqueueTask(&CustomTask{
		Id:        uint(111111),
		Data:      rand.Intn(1000),
	})
	wp.EnqueueTask(&CustomTask{
		Id:        uint(123124),
		Data:      rand.Intn(1000),
	})

	wp.Start()

	for i := 1; i <= 1000; i++ {
		t := CustomTask{
			Id:        uint(i),
			Data:      rand.Intn(1000),
		}
		wp.EnqueueTask(&t)
	}

}
