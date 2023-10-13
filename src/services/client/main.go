package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/kodmain/kitsune/src/internal/core/server/protocols/tcp"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/shirou/gopsutil/cpu"
)

var max = 10000000
var total = 0
var rps = 0
var totalTime time.Duration // Variable to keep track of total time
var mu sync.Mutex
var worker = 1

func main() {
	run("9999")

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		percent, err := cpu.Percent(0, false)
		if err != nil {
			log.Println("Erreur lors de la récupération de l'utilisation du CPU:", err)
		}

		cpuUsage := 0.0
		if len(percent) > 0 {
			cpuUsage = percent[0]
		}

		avgTime := float64(totalTime.Milliseconds()) / float64(total) // Calculate average time in ms

		log.Printf("Count Request/Sec: %d, REQS: %d/%d, Go Routine: %d, MemoryUsage: %d Mb, CPU Usage: %.2f%%, Avg Time per Request: %.2fms", rps, total, max, runtime.NumGoroutine(), bToMb(m.Alloc), cpuUsage, avgTime)
		if int(total) >= max {
			ticker.Stop()
			break
		}

		mu.Lock()
		rps = 0
		mu.Unlock()
	}
}

func run(port string) {
	client := tcp.NewClient() // youka-PRODUCTION-9de5d4b457bad9c7.elb.eu-west-3.amazonaws.com
	service1, err := client.Connect("127.0.0.1", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for j := 0; j < worker; j++ {
		go func(j int) {
			for i := 0; i < max; i++ {
				query1 := service1.MakeExchange()
				client.Send(func(responses ...*transport.Response) {
					if query1.Request().Id == responses[0].Id {
						mu.Lock()
						rps++
						total++
						mu.Unlock()
					}
				}, query1)
			}
		}(j)
	}
}

/*
func run(port string) {
	client := socket.NewClient("127.0.0.1:" + port) // youka-PRODUCTION-9de5d4b457bad9c7.elb.eu-west-3.amazonaws.com
	err := client.Connect()

	if logger.Error(err) {
		os.Exit(1)
	}

	for j := 0; j < worker; j++ {
		go func(j int) {
			for i := 0; i < max; i++ {
				req := transport.CreateRequest()

				start := time.Now() // Record the start time
				client.SendSync(req)
				elapsed := time.Since(start) // Calculate elapsed time

				mu.Lock()
				total++
				rps++
				totalTime += elapsed // Update total time
				mu.Unlock()
			}
		}(j)
	}
}
*/

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
