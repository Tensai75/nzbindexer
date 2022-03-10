package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var counter int

var wg sync.WaitGroup

func main() {

	counter = 0
	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()

	guard := make(chan struct{}, conf.ParallelScans)

	running := true

	// endlessly loop until process is aborted (e.g. via Ctrl-C)
	// the program will end gracefully when aborted
	for running {
		select {
		case <-ctx.Done():
			running = false
		default:
			for _, group := range serverGroups {
				select {
				case <-ctx.Done():
					break
				default:
					guard <- struct{}{} // will block if guard channel is already filled
					go func(group string) {
						select {
						case <-ctx.Done():
							return
						default:
							wg.Add(1)
							indexer(group, ctx)
							<-guard
						}
					}(group)
				}
			}
		}

	}

	wg.Wait()
	duration := time.Since(start)
	perSecond := float64(counter) / duration.Seconds()
	fmt.Printf("A total of %d messages were processed in %v (%d Messages/s)\n", counter, duration, int(perSecond))

}

func init() {

	fmt.Println("Loading configuration")
	if err := loadConfig(); err != nil {
		fmt.Println("Fatal error while loading configuration file!")
		os.Exit(1)
	}

	fmt.Println("Connecting to database")
	if err := connectMySQL(); err != nil {
		fmt.Println("Fatal error while connecting to database!")
		os.Exit(1)
	}

	fmt.Println("Scanning the groups")
	if err := scanGroups(); err != nil {
		fmt.Println("Fatal error while scanning the groups!")
		os.Exit(1)
	}
}
