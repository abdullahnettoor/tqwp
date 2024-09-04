package main

import (
	"fmt"
	"log"

	"math/rand"
)

type Task struct {
	Id         uint
	Data       any
	Retries    int
	MaxRetries int
}

func (t *Task) Process() error {
	log.Printf("--- Processing TASK %d: Retries %d", t.Id, t.Retries)
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
