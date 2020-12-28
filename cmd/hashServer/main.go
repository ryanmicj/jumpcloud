package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	shutdownChannel := make(chan string)

	listenPort := 8080

	//Check to see if a port was supplied on the command line, and that it is numeric
	var listenPortString string
	if len(os.Args) > 1 {
		listenPortString = os.Args[1]
		if len(listenPortString) > 0 {
			portArg, err := strconv.Atoi(listenPortString)
			if err == nil {
				listenPort = portArg
			}
		}
	}
	listenPortString = strconv.Itoa(listenPort)

	fmt.Println("Using port " + listenPortString)
	go listen(listenPortString, shutdownChannel, &wg)

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
