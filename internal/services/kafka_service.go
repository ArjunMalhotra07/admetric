package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ArjunMalhotra/internal/model"
	"github.com/Shopify/sarama"
)

const (
	clickTopic = "ad-clicks"
)

type KafkaService struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	config   *sarama.Config
	brokers  []string
}

func NewKafkaService(brokers []string) (*KafkaService, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	return &KafkaService{
		producer: producer,
		consumer: consumer,
		config:   config,
		brokers:  brokers,
	}, nil
}

func (s *KafkaService) PublishClick(click model.Click) error {
	msg, err := json.Marshal(click)
	if err != nil {
		return fmt.Errorf("failed to marshal click: %v", err)
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: clickTopic,
		Value: sarama.StringEncoder(msg),
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (s *KafkaService) StartConsumer(clickService *ClickService) error {
	partitionConsumer, err := s.consumer.ConsumePartition(clickTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("failed to create partition consumer: %v", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				var click model.Click
				if err := json.Unmarshal(msg.Value, &click); err != nil {
					log.Printf("Failed to unmarshal click: %v", err)
					continue
				}

				if err := clickService.ProcessClick(click); err != nil {
					log.Printf("Failed to process click: %v", err)
				}
			case <-signals:
				return
			}
		}
	}()

	return nil
}

func (s *KafkaService) Close() error {
	if err := s.producer.Close(); err != nil {
		return fmt.Errorf("failed to close producer: %v", err)
	}
	if err := s.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %v", err)
	}
	return nil
}
