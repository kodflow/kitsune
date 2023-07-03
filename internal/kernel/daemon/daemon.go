package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kodmain/KitsuneFramework/internal/env"
)

type Handler struct {
	Name string
	Call func() error
}

var sigs chan os.Signal = make(chan os.Signal, 1)
var done chan bool = make(chan bool, 1)

func Start(handlers ...*Handler) {
	fmt.Println(env.BUILD_APP_NAME, "start")

	if _, err := GetPID(env.BUILD_APP_NAME); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := SetPID(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		ClearPID(env.BUILD_APP_NAME)
		done <- true
	}()

	for _, handler := range handlers {
		go func(handler *Handler) {
			count := 0
			startTime := time.Now()

			for {
				fmt.Println(env.BUILD_APP_NAME, handler.Name, "start")
				if err := handler.Call(); err != nil {
					fmt.Println(env.BUILD_APP_NAME, handler.Name, "fail")
					if count >= 2 {
						elapsedTime := time.Since(startTime)
						if elapsedTime < time.Minute {
							fmt.Println(env.BUILD_APP_NAME, handler.Name, "exit")
							done <- true
							break
						}
						count = 0
						startTime = time.Now()
					}
					count++
				} else {
					break
				}
			}
		}(handler)
	}

	<-done
	fmt.Println(env.BUILD_APP_NAME, "exit")
}

func Stop() {
	fmt.Println(env.BUILD_APP_NAME, "stop")
	sigs <- syscall.SIGTERM
}
