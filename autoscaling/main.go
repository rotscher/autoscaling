package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

const meterName = "github.com/rotscher/autoscaling/autoscaling"
const coresCount = 1
const percentage = 20
const timeSeconds = 600

type queue struct {
	count int64
}

type basicAuth struct {
	user     string
	password string
	enabled  bool
}

func main() {

	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	auth := basicAuth{enabled: false}

	if len(user) > 0 && len(password) > 0 {
		auth.user = user
		auth.password = password
		auth.enabled = true
	}

	fmt.Printf("Starting autoscaling demo app with following params: coresCount=%d, percentage=%d, timeSecondes=%d", coresCount, percentage, timeSeconds)
	ctx := context.Background()
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(meterName)

	queue := queue{count: 0}
	queueCount, err := meter.Int64Gauge("queue_current_count", api.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}

	//init the metric
	queueCount.Record(ctx, queue.count)

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/cpu", func(w http.ResponseWriter, r *http.Request) {
		if auth.enabled == true {
			username, p, ok := r.BasicAuth()
			if !ok || username != auth.user || p != auth.password {
				w.WriteHeader(401)
				return
			}
		}

		_, _ = io.WriteString(w,
			fmt.Sprintf("this generates some cpu load: coresCount=%d, percentage=%d, timeSecondes=%d\n", coresCount, percentage, timeSeconds))
		fmt.Printf("got /cpu request\n")
		runCpu()
	})
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		queue.count = queue.count + 1
		queueCount.Record(ctx, queue.count)
	})

	http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		if queue.count > 0 {
			queue.count = queue.count - 1
			queueCount.Record(ctx, queue.count)
		}
	})

	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	_ = http.ListenAndServe(":3333", nil)
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	_, _ = io.WriteString(w, "This is my website!\n")
}
func runCpu() {

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
