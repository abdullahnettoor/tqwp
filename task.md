Build a Concurrent Task Queue with Retry Mechanism
--------------------------------------------------

**Objective**:

Create a Go application that simulates a task queue system. This system will handle tasks concurrently, implement a retry mechanism for failed tasks, and log results. The task queue should be able to process tasks from an in-memory queue and retry failed tasks a specified number of times before giving up.

Requirements:
-------------

**Task Definition**

Define a task as a struct with at least the following fields:

*   ID: A unique identifier for the task.
    
*   Data: Some data to be processed (you can use a string or struct for this).
    
*   Retries: The number of times the task has been retried.

---

**Task Queue:**

*   Implement an in-memory task queue to store and manage tasks. Use Go channels to simulate task enqueuing and dequeuing.
    
---

**Task Processor:**

*   Create a worker pool that processes tasks concurrently. Each worker should process tasks from the queue.
    
*   Implement a retry mechanism where tasks that fail are retried up to a specified number of times. After reaching the maximum number of retries, the task should be logged as failed.
    
---

**Error Handling and Logging:**

*   Simulate task processing with possible failures. For example, randomly fail a task based on some probability.
    
*   Implement logging for task processing, retries, and failures. Use Go’s built-in logging package or a third-party logger.
    
---
**Concurrency**:

*   Use goroutines to handle multiple workers processing tasks concurrently.
    
*   Ensure that the system can handle a high volume of tasks efficiently.
    
---

**Task Completion and Statistics:**

*   Track and log the number of tasks processed, successfully completed tasks, and failed tasks.
    
*   Print or log summary statistics after all tasks have been processed.

---
Happy Coding 😊