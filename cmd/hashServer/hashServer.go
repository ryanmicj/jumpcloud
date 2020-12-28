package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Instantiate the hashServer and begin listening on the provided port.
func listen(port string, shutdownChannel chan string, wg *sync.WaitGroup) {

	/*
	 We will use a Reader/Writer Mutex to prevent concurrent updated to
	 the current request counter and map of hashed passwords.  This will allow
	 multiple GET requests to proceed without blocking, which is beneficial,
	 since they will be unlikely to accessing the same index
	*/
	count := 0
	hashedPasswords := make(map[int]string)
	var mutex = &sync.RWMutex{}

	/*
	 We will use a closure to update the statistics for our hash server.
	 This closure will initiate a go routine update the statistics, as well as
	 provide the handler method to respond to web requests.
	*/
	statsChannel := make(chan int64)
	statsHandler := getStatsHandler(statsChannel)

	//Instantiate a handler to accept the hash request and generate an ID for that hash
	//The actual hashing will happen asynchronously in a go routine
	hashHandler := func(w http.ResponseWriter, req *http.Request) {
		if checkMethodType(req, "POST") == true {
			//Fetch the password from the form payload
			passwordToHash := req.FormValue("password")
			if len(passwordToHash) == 0 {
				log.Printf("Error Reading Form Data: \"password\" not found")
				http.Error(w, "Error Reading Request Body", http.StatusBadRequest)
			} else {
				mutex.Lock()
				count++
				current := count
				mutex.Unlock()

				io.WriteString(w, fmt.Sprintf("%d", current))
				wg.Add(1)

				//Spin up a go routine to do the actual hashing
				go func(passwordToHash string, index int) {
					defer wg.Done()

					startTime := time.Now()
					hashedPassword := encode(passwordToHash, fetchSha512Hash())
					endTime := time.Now()

					//Sleep for 5 seconds before the hashed password is available to GET
					time.Sleep(5 * time.Second)

					mutex.Lock()
					hashedPasswords[index] = hashedPassword
					mutex.Unlock()

					//Inform the stats channel how long this request took
					requestLength := endTime.Sub((startTime)).Microseconds()
					statsChannel <- requestLength
				}(passwordToHash, current)
			}
		} else {
			writeInvalidMethodError(w, req.Method)
		}
	}

	//Instantiate a handler to handle GET requests for each hashId
	fetchHashHandler := func(w http.ResponseWriter, req *http.Request) {
		if checkMethodType(req, "GET") == true {
			//Split the request URL into 3 parts:
			// 1. "/"
			// 2. "hash" (or whatever the context root is assigned)
			// 3. the id of the ahs to fetch
			//Requests of any other form will return an error
			requestPath := strings.SplitAfter(req.URL.String(), "/")
			if len(requestPath) != 3 {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
			} else {
				hashRequest, err := strconv.Atoi(requestPath[2])
				if err != nil {
					http.Error(w, "Invalid Request - must be numeric: "+"/"+requestPath[2], http.StatusBadRequest)
				} else {
					//Look up this hashId in our Map.
					mutex.RLock()
					hashedPassword, found := hashedPasswords[hashRequest]
					mutex.RUnlock()

					if found == true && len(hashedPassword) > 0 {
						io.WriteString(w, hashedPassword)
					} else {
						http.Error(w, "Not Found", http.StatusNotFound)
					}
				}
			}
		} else {
			writeInvalidMethodError(w, req.Method)
		}
	}

	shutdownHandler := func(w http.ResponseWriter, req *http.Request) {
		if checkMethodType(req, "POST") == true {
			shutdownChannel <- "shutdown"
			close(shutdownChannel)
		} else {
			writeInvalidMethodError(w, req.Method)
		}
	}

	http.HandleFunc("/hash/", fetchHashHandler)    // Handle GET requests: GET /hash/{hashId}
	http.HandleFunc("/hash", hashHandler)          // Handle POST Requests: POST /hash + Form payload
	http.HandleFunc("/shutdown", shutdownHandler)  // Handle *any* request: GET/POST /shutdown
	http.HandleFunc(("/statistics"), statsHandler) // Handle GET requests: GET /statistics

	//TODO: Determine how to gracefully shutdown the server.
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

//Check that the Request method matches what the handler expects/allows
func checkMethodType(r *http.Request, allowedMethod string) bool {
	actualMethod := r.Method
	if strings.Compare(actualMethod, allowedMethod) == 0 {
		return true
	}

	return false
}

//Write the standard error Response for a disallowed method
func writeInvalidMethodError(w http.ResponseWriter, method string) {
	http.Error(w, "Invalid Method: "+method, http.StatusMethodNotAllowed)
}
