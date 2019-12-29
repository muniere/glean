package main

import (
	"log"

	"github.com/muniere/glean/internal/app/server"
	"github.com/muniere/glean/internal/pkg/rpc"
	"github.com/muniere/glean/internal/pkg/task"
)

func main() {
	queue := task.NewQueue()

	consumer := server.NewConsumer(queue, server.ConsumerConfig{
		Concurrency: task.Concurrency,
	})

	producer := server.NewProducer(queue, server.ProducerConfig{
		Address: rpc.LocalAddr,
		Port:    rpc.Port,
	})

	var err error

	err = consumer.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = producer.Start()
	if err != nil {
		log.Fatal(err)
	}

	consumer.Wait()
	producer.Wait()
}
