package goroutines

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup, jobs <-chan string) {
	defer wg.Done()
	if len(jobs) > 0 {
		fmt.Printf("worker:%d spawning\n", id)
		for j := range jobs {
			task_time, _ := time.ParseDuration(j + "s")
			fmt.Printf("worker:%d sleep:%s\n", id, j)
			time.Sleep(task_time)
		}
		fmt.Printf("worker:%d stopping\n", id)
	}
}

func Run(poolSize int) {
	jobs := make(chan string, poolSize)
	var wg sync.WaitGroup
	id := 1
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str := scanner.Text()
		jobs <- str
		if id <= poolSize {
			wg.Add(1)
			go worker(id, &wg, jobs)
			id++
		}
	}
	close(jobs)
	wg.Wait()
}
