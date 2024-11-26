package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

const meterName = "github.com/rotscher/autoscaling/autoscaling"

func main() {
	ctx := context.Background()
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(meterName)

	// Start the prometheus HTTP server and pass the exporter Collector to it
	//go serveMetrics()

	queueCount, err := meter.Int64UpDownCounter("queue_current_count", api.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/cpu", getCPU)
	http.HandleFunc("/memory", getMemory)
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		queueCount.Add(ctx, 1)
	})

	http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		queueCount.Add(ctx, -1)
	})

	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

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
