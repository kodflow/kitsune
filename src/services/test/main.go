package main

import (
	"flag"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/kodmain/kitsune/src/internal/core/server/protocols/tcp"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/shirou/gopsutil/cpu"
)

var (
	URL    = flag.String("url", "localhost", "set url")
	PORT   = flag.String("port", "9999", "set port")
	NUMCPU = flag.Int("cpu", 1, "set max CPU")
)

func init() {
	if !strings.HasSuffix(os.Args[0], ".test") {
		flag.Parse()
	}
}

func main() { //runtime.NumCPU()
	runtime.GOMAXPROCS(*NUMCPU) // facultatif

	var mu sync.Mutex
	var max = 10000000 / *NUMCPU // divisez par le nombre de CPU pour garder le même total
	var rps = 0
	var total = 0

	var rpsHistory []int
	var minRPS int = math.MaxInt64
	var maxRPS int = math.MinInt64
	var sumRPS int = 0

	for i := 0; i < *NUMCPU; i++ {
		go func() {
			client := tcp.NewClient()
			service, _ := client.Connect(*URL, *PORT)
			for i := 0; i < max; i++ {
				query1 := service.MakeExchange()
				client.Send(func(responses ...*transport.Response) {
					mu.Lock()
					rps++
					total++
					mu.Unlock()
				}, query1)
			}
		}()
	}

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		percent, _ := cpu.Percent(0, false)
		cpuUsage := 0.0
		if len(percent) > 0 {
			cpuUsage = percent[0]
		}

		rpsHistory = append(rpsHistory, rps)
		if rps < minRPS {
			minRPS = rps
		}
		if rps > maxRPS {
			maxRPS = rps
		}
		sumRPS += rps

		avgRPS := sumRPS / len(rpsHistory)
		deltaRPS := maxRPS - minRPS

		log.Printf("Request/Sec: %d, Avg: %d, Min: %d, Max: %d, Delta: %d, REQS: %d/%d, Go Routine: %d, MemoryUsage: %d Mb, CPU Usage: %.2f%%", rps, avgRPS, minRPS, maxRPS, deltaRPS, total, max, runtime.NumGoroutine(), bToMb(m.Alloc), cpuUsage)
		rps = 0
		if int(total) >= max**NUMCPU {
			ticker.Stop()
			break
		}
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
