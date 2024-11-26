package main

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"
)

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/cpu", getCPU)
	http.HandleFunc("/memory", getMemory)

	_ = http.ListenAndServe(":3333", nil)
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	_, _ = io.WriteString(w, "This is my website!\n")
}
func getCPU(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /cpu request\n")
	_, _ = io.WriteString(w, "this generates some cpu load!\n")

	coresCount := 1
	percentage := 50
	timeSeconds := 120
	mainBegin := time.Now()
	runtime.GOMAXPROCS(coresCount)

	// second     ,s  * 1
	// millisecond,ms * 1000
	// microsecond,μs * 1000 * 1000
	// nanosecond ,ns * 1000 * 1000 * 1000

	// every loop : run + sleep = 1 unit

	// 1 unit = 100 ms may be the best
	unitHundredsOfMicrosecond := 1000
	runMicrosecond := unitHundredsOfMicrosecond * percentage
	sleepMicrosecond := unitHundredsOfMicrosecond*100 - runMicrosecond
	for i := 0; i < coresCount; i++ {
		go func(id int) {
			runtime.LockOSThread()
			// endless loop
			for {
				begin := time.Now()
				for {
					// run 100%
					if time.Now().Sub(begin) > time.Duration(runMicrosecond)*time.Microsecond {
						break
					}
				}
				// sleep
				time.Sleep(time.Duration(sleepMicrosecond) * time.Microsecond)
				if time.Now().Sub(mainBegin) > time.Duration(timeSeconds)*time.Second {
					fmt.Printf("finished.... %d, %d, %d\n", id, time.Now().Sub(mainBegin), time.Duration(timeSeconds)*time.Second)
					break
				}
			}
		}(i)
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func getMemory(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /memory request\n")
	_, _ = io.WriteString(w, "this should use some memory!\n")
}
