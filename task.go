package tqwp

import (
	"fmt"

	"math/rand"
)

type Task struct {
	Id      uint
	Data    any
	Retries int
}

func (t *Task) Process() error {
	num, isInt := t.Data.(int)
	if !isInt {
		return fmt.Errorf("invalid type")
	}
	divisor := rand.Intn(2)
	if divisor == 0 {
		return fmt.Errorf("division by zero")
	}
	t.Data = num / divisor
	return nil
}
