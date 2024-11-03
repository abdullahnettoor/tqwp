/*
Package tqwp implements a concurrent task queue worker pool with retry mechanism.

Basic usage:

	wp := tqwp.New(&tqwp.WorkerPoolConfig{
	    NumOfWorkers: 10,
	    MaxRetries:   3,
	    QueueSize:    100,
	})
	defer wp.Stop()

	wp.Start()
	wp.EnqueueTask(&CustomTask{})

For more examples, see: https://github.com/abdullahnettoor/tqwp/tree/main/examples
*/
package tqwp