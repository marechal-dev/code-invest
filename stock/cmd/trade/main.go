package main

import (
	"encoding/json"
	"fmt"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/marechal-dev/code-invest/stock/internal/domain/market/dtos"
	"github.com/marechal-dev/code-invest/stock/internal/domain/market/entities"
	"github.com/marechal-dev/code-invest/stock/internal/domain/market/transformers"
	"github.com/marechal-dev/code-invest/stock/internal/infra/messaging/kafka"
)

func main() {
	ordersInputChannel := make(chan *entities.Order)
	ordersOutputChannel := make(chan *entities.Order)

	waitGroup := &sync.WaitGroup{}
	defer waitGroup.Wait()

	kafkaMessageChannel := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "latest",
	}
	producer := kafka.NewProducer(configMap)
	consumer := kafka.NewConsumer(configMap, []string{"input"})

	go consumer.Consume(kafkaMessageChannel)

	book := entities.NewBook(ordersInputChannel, ordersOutputChannel, waitGroup)
	go book.Trade()

	go func() {
		for message := range kafkaMessageChannel {
			waitGroup.Add(1)

			fmt.Println(string(message.Value))

			tradeInput := dtos.TradeInput{}

			err := json.Unmarshal(message.Value, &tradeInput)
			if err != nil {
				panic(err)
			}

			order := transformers.TransformInput(tradeInput)
			ordersInputChannel <- order
		}
	}()

	for response := range ordersOutputChannel {
		output := transformers.TransformOutput(response)
		outputJson, err := json.MarshalIndent(output, "", "   ")
		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}

		producer.Publish(outputJson, []byte("orders"), "output")
	}
}
