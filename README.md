# periodic

periodic is a Golang task scheduling pkg. This helps you to run functions periodically using pre-determined interval easily.
> I think this pkg is not stable. So, I would not recommend using it in prod env.
## Quickstart

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blackironj/periodic"
)

func helloWorld() {
	fmt.Println("Hello World")
}

func helloWithParams(first, second string) {
	fmt.Printf("Hello, %s %s \n", first, second)
}

func main() {
	scheduler := periodic.NewScheduler()

	//Register tasks
	helloWorld, _ := periodic.NewTask(helloWorld)
	scheduler.RegisterTask("task1", time.Millisecond*100, helloWorld)

	helloWithParams, _ := periodic.NewTask(helloWithParams, "alice", "bob")
	scheduler.RegisterTask("task2", time.Second*2, helloWithParams)

	//Run tasks
	fmt.Println("Run tasks")
	scheduler.Run()

	//Stop tasks before program is shutting down
	defer func() {
		scheduler.Stop()
		fmt.Println("Every task is stopped")
	}()

	//Expected output:
	//Initially, the task starts right away, and then next tick depends on interval time
	//Print "Hello world" every 100 milliseconds
	//Print "Hello, alice, bob" every 2 seconds

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
}
```
> The shorter the interval, the harder the task is to be excuted at the right time.