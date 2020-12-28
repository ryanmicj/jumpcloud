package main

import (
	"fmt"
	"testing"
)

func TestComputeStats(t *testing.T) {
	var mockStats Statistics
	var mockDuration int64

	//Update our stats 5 times
	var updateCount int = 5
	var averageDuration int64 = 100
	var averageDurationFloat float64 = 100.0
	for i := 0; i < updateCount; i++ {
		mockDuration = updateStats(&mockStats, mockDuration, averageDuration)
	}

	if mockStats.Total != updateCount {
		msg := fmt.Sprintf("incorrect number of requests counted. Expected %d but counted %d", updateCount, mockStats.Total)
		t.Fatal(msg)
	}

	if mockStats.Average != averageDurationFloat {
		msg := fmt.Sprintf("incorrect average for requests. Expected %f but computed %f", averageDurationFloat, mockStats.Average)
		t.Fatal(msg)
	}

}
