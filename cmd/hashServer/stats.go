package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

/*
Statistics type definition for the Statistics type that will be returned when
making calls to the "/statistics" endpoint
*/
type Statistics struct {
	Total   int     `json:"total"`
	Average float64 `json:"average"`
}

//Get the Http Handler function that will respond to statisics requests.
//This method will start a go routine that will compute the average time
func getStatsHandler(statsChan chan int64) func(w http.ResponseWriter, req *http.Request) {
	var requestStats Statistics
	var totalDuration int64

	//We will use a Mutex to synchrnize access to the statistics data.  Otherwise, it
	//would be possible to fetch the data while it is in the process of being updated.
	//We will use a basic mutext since the number of writes will likely exceed the
	//number of reads of the statistics data
	var mutex = &sync.Mutex{}

	//Run a go routing that reads from the statsChannel and updates the statistics for this channel
	go func(statsChan chan int64) {
		for requestDuration := range statsChan {

			mutex.Lock()
			totalDuration = updateStats(&requestStats, totalDuration, requestDuration)
			mutex.Unlock()
		}
	}(statsChan)

	return func(w http.ResponseWriter, req *http.Request) {
		if checkMethodType(req, "GET") {
			mutex.Lock()
			output, err := json.Marshal(requestStats)
			mutex.Unlock()
			if err != nil {
				fmt.Println("Error writing statistics output to JSON!")
			} else {
				fmt.Println(string(output))
				io.WriteString(w, string(output))
			}
		} else {
			writeInvalidMethodError(w, req.Method)
		}
	}
}

func updateStats(stats *Statistics, totalDuration int64, currentDuration int64) int64 {
	stats.Total++
	totalDuration += currentDuration
	stats.Average = (float64(totalDuration)) / (float64(stats.Total))

	return totalDuration
}
