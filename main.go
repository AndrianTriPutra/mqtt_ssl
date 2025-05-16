package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mqTT/logger"
	"mqTT/mqtt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func init() {
	log.Println("=================================")
	log.Println("app    :mqtt_ssl")
	log.Println("author :@tp")
	log.Println("release:2025-04-04T00:00")
	log.Println("=================================")
}

func main() {
	logger.Load("debug")
	ctx := context.Background()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	pathca := fmt.Sprintf("%s/docker/ssl/ca.pem", basepath)

	topic := "ssl/seapoint/"
	setting := mqtt.Setting{
		Path: pathca,

		ServerName: "127.0.0.1", //alt_names
		Host:       "ssl://localhost:8883",
		User:       "thai",
		Pass:       "chan",

		ClientID: fmt.Sprintf("mqtt_tls/%d", time.Now().Unix()),
		Qos:      1,
		Topic:    topic + "#",
	}

	repomqtt := mqtt.NewMQTT(setting)
	err := repomqtt.Connect(ctx)
	if err != nil {
		logger.Level("fatal", "main", "error on connect mqtt:"+err.Error())
	}

	stoped := make(chan os.Signal, 1)
	signal.Notify(stoped,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	type Msg struct {
		Sender  string `json:"sender"`
		Message string `json:"message"`
	}
	counter := uint(0)
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		defer wg.Done()
		runtime.Gosched()
		for {
			<-ticker.C
			counter++
			message := Msg{
				Sender:  "client-1",
				Message: fmt.Sprintf("Hello form SSL-%d", counter),
			}

			js, _ := json.Marshal(message)
			logger.Level("debug", "main", "Send:"+string(js))
			err := repomqtt.Publish(ctx, topic+message.Sender, string(js))
			if err != nil {
				logger.Level("error", "main", "error Publish mqtt:"+err.Error())
			}
		}
	}()

	go func() {
		defer wg.Done()
		runtime.Gosched()
		for {
			msg := <-repomqtt.OnMsg()
			logger.Level("debug", "main", "\t\t Recv:"+string(msg))
			logger.Level("debug", "main", "")
		}
	}()

	message := ""
	s := <-stoped
	switch s {
	case syscall.SIGHUP:
		message = "[hungup]"
	case syscall.SIGINT:
		message = "[interupt]"
	case syscall.SIGTERM:
		message = "[force stop]"
	case syscall.SIGQUIT:
		message = "[stop and core dump]"
	default:
		message = "[unknown signal]"
	}
	logger.Level("info", "Run", message)
}
