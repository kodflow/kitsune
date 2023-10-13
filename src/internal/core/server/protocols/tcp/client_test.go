package tcp_test

import (
	"log"
	"math"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/kodmain/kitsune/src/internal/core/server/protocols/tcp"
	"github.com/kodmain/kitsune/src/internal/core/server/transport"
	"github.com/shirou/gopsutil/cpu"
)

func TestClient(t *testing.T) {
	server1 := tcp.NewServer("localhost:8080")
	server1.Start()
	defer server1.Stop()

	client := tcp.NewClient()
	service1, _ := client.Connect("localhost", "8080")
	query1 := service1.MakeExchange()
	response := make(chan bool)

	client.Send(func(responses ...*transport.Response) {
		response <- true
	}, query1)
	<-response
}

func BenchmarkLocal(b *testing.B) {
	server1 := tcp.NewServer("localhost:8080")
	server1.Start()
	defer server1.Stop()

	client := tcp.NewClient()
	service1, _ := client.Connect("localhost", "8080")

	b.Run("benchmark", func(b *testing.B) {
		b.ResetTimer() // Ne compte pas la configuration initiale dans le temps de benchmark
		var max = 1000000
		var rps = 0
		var total = 0
		var mu sync.Mutex

		var rpsHistory []int
		var minRPS int = math.MaxInt64
		var maxRPS int = math.MinInt64
		var sumRPS int = 0

		go func() {
			for i := 0; i < max; i++ {
				query1 := service1.MakeExchange(false)
				client.Send(func(responses ...*transport.Response) {
					mu.Lock()
					rps++
					total++
					mu.Unlock()
				}, query1)
			}
		}()

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
			if int(total) >= max {
				ticker.Stop()
				break
			}
		}
	})
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

/*



	var address, port, protocol string = "", "", ""




func TestClient(t *testing.T) {
	logger.SetLevel(levels.OFF)
	server := tcp.NewServer("localhost:8080")
	server.Start()
	defer server.Stop()

	client := tcp.NewClient("localhost:8080")

	t.Run("NewClient", func(t *testing.T) {
		assert.NotNil(t, client, "Expected client to be non-nil")
		assert.Equal(t, "localhost:8080", client.Address)
	})

	t.Run("ConnectAndClose", func(t *testing.T) {
		err := client.Connect()
		assert.Nil(t, err, "Expected no error on connecting")

		err = client.Connect()
		assert.NotNil(t, err, "Expected error on connecting already done")

		err = client.Disconnect()
		assert.Nil(t, err, "Expected no error on closing")

		err = client.Disconnect()
		assert.NotNil(t, err, "Expected error on closing already closed connection")

		err = client.Connect()
		assert.Nil(t, err, "Expected no error on connecting")
	})

	t.Run("Max", func(t *testing.T) {
		const max = 1000
		clients := map[int]*tcp.Client{}
		for c := 0; c < max; c++ {
			clients[c] = tcp.NewClient("localhost:8080")
			err := clients[c].Connect()
			assert.Nil(t, err, "Expected no error on multiple connect")
		}

		for c := 0; c < max; c++ {
			err := clients[c].Disconnect()
			assert.Nil(t, err, "Expected no error on multiple connect")
		}
	})

	t.Run("RequestAndResponse", func(t *testing.T) {
		client.Connect()
		req := transport.CreateRequest()
		promise, err := client.Send(req)
		assert.Nil(t, err, "Expected no error on sending request")
		assert.NotNil(t, promise, "Expected promise of response")
		res := promise.Wait()
		assert.NotNil(t, res, "Expected response")
		assert.Equal(t, res.Id, req.Id, "Expected promise of response")
	})

	t.Run("RequestOnly", func(t *testing.T) {
		client.Connect()
		req := transport.CreateRequestOnly()
		promise, err := client.Send(req)
		assert.Nil(t, err, "Expected no error on sending request")
		assert.NotNil(t, promise, "Expected promise of response")
	})

}

func BenchmarkRequestsOnly(b *testing.B) {
	logger.SetLevel(levels.OFF)
	server := tcp.NewServer("localhost:8080")
	server.Start()

	client := tcp.NewClient("localhost:8080")
	client.Connect()

	b.Run("benchmark", func(b *testing.B) {
		b.ResetTimer() // Ne compte pas la configuration initiale dans le temps de benchmark
		var max = 1000000
		var worker = runtime.NumCPU()
		var total int32
		var errors int32
		var rps int32
		var rpw = max / worker
		var mu sync.Mutex

		for j := 0; j < worker; j++ {
			go func(j int) {
				var i int
				for i = 0; i < rpw; i++ {
					req := transport.CreateRequestOnly()
					_, err := client.Send(req)
					if err != nil {
						log.Println(err)
						mu.Lock()
						errors++
						mu.Unlock()
						continue
					}

					mu.Lock()
					total++
					rps++
					mu.Unlock()
				}

				log.Printf("Finished worker %d after %d requests", j, i)
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
			rps = 0
			if int(total) >= max {
				ticker.Stop()
				break
			}
		}
	})

	client.Disconnect()
	server.Stop()
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func BenchmarkRequestAndResponse(b *testing.B) {
	logger.SetLevel(levels.OFF)

	server := tcp.NewServer("localhost:8080")
	server.Start()

	client := tcp.NewClient("localhost:8080")
	client.Connect()

	b.Run("benchmark", func(b *testing.B) {
		b.ResetTimer() // Ne compte pas la configuration initiale dans le temps de benchmark
		var max = 1000000
		var worker = 1
		var total int32
		var errors int32
		var rps int32
		var rpw = max / worker
		var mu sync.Mutex

		for j := 0; j < worker; j++ {
			go func(j int) {
				var i int
				for i = 0; i < rpw; i++ {
					req := transport.CreateRequest()
					_, err := client.SendSync(req)
					if err != nil {
						log.Println(err)
						mu.Lock()
						errors++
						mu.Unlock()
						continue
					}

					mu.Lock()
					total++
					rps++
					mu.Unlock()
				}

				log.Printf("Finished worker %d after %d requests", j, i)
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
			rps = 0
			if int(total) >= max {
				ticker.Stop()
				break
			}
		}
	})

	client.Disconnect()
	server.Stop()
}

/*
func BenchmarkLocal(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	client := tcp.NewClient("youka-PRODUCTION-9de5d4b457bad9c7.elb.eu-west-3.amazonaws.com:9999")
	client.Connect()
	time.Sleep(time.Second)

	b.Run("benchmark", func(b *testing.B) {
		b.ResetTimer() // Ne compte pas la configuration initiale dans le temps de benchmark
		var max = 10000000
		var worker = runtime.NumCPU() * 30
		var total int32
		var errors int32
		var rps int32
		var rpw = max / worker
		var mu sync.Mutex

		for j := 0; j < worker; j++ {
			go func(j int) {
				var i int
				for i = 0; i < rpw; i++ {
					req := transport.CreateRequest()
					res, err := client.SendSync(req)
					if err != nil {
						log.Println(err)
						mu.Lock()
						errors++
						mu.Unlock()
						continue
					}

					if res.Id != req.Id {
						log.Println("Mismatch ID")
						mu.Lock()
						errors++
						mu.Unlock()
					}

					mu.Lock()
					total++
					rps++
					mu.Unlock()
				}

				log.Printf("Finished worker %d after %d requests", j, i)
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
			rps = 0
			if int(total) >= max {
				ticker.Stop()
				break
			}
		}
	})

	client.Disconnect()
}
*/
