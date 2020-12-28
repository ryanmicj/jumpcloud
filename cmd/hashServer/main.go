package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	shutdownChannel := make(chan string)

	fmt.Println("Using port " + os.Args[1])
	go listen(os.Args[1], shutdownChannel, &wg)

	//Wait for the shutdown message, then exit
	for msg := range shutdownChannel {
		if strings.Compare(msg, "shutdown") == 0 {
			fmt.Println("Shutdown Channel closed - shutting down main thread.")
		}
	}

	//Need to perform a graceful shutdown of the http server thread.

	//Wait to ensure that any workers have completed before exiting
	wg.Wait()

	fmt.Println("Exiting.")
}
