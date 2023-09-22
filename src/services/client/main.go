package main

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/kodmain/kitsune/src/internal/core/server/protocols/socket"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/shirou/gopsutil/cpu"
)

func main() {
	var max = 100000
	var total = 0
	var rps = 0
	var mu sync.Mutex
	var worker = runtime.NumCPU()

	client := socket.NewClient("youka-PRODUCTION-9de5d4b457bad9c7.elb.eu-west-3.amazonaws.com:9999")
	client.Connect()

	for j := 0; j < worker; j++ {
		go func(j int) {
			for i := 0; i < max; i++ {
				req := transport.CreateRequest()
				client.SendSync(req)
				mu.Lock()
				total++
				rps++
				mu.Unlock()
			}
		}(j)
	}

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

		log.Printf("Request/Sec: %d, REQS: %d/%d, Go Routine: %d, MemoryUsage: %d Mb, CPU Usage: %.2f%%", rps, total, max, runtime.NumGoroutine(), bToMb(m.Alloc), cpuUsage)
		if int(total) >= max {
			ticker.Stop()
			break
		}

		mu.Lock()
		rps = 0
		mu.Unlock()
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
