package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ArjunMalhotra/internal/model"
	"github.com/Shopify/sarama"
)

const (
	clickTopic = "ad-clicks"
	maxRetries = 5
	retryDelay = 2 * time.Second
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
	config.Producer.Retry.Max = maxRetries
	config.Producer.Retry.Backoff = retryDelay

	// Add timeout settings
	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second

	// Try to connect with retries
	var producer sarama.SyncProducer
	var err error

	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to connect to Kafka brokers: %v (attempt %d/%d)", brokers, i+1, maxRetries)
		producer, err = sarama.NewSyncProducer(brokers, config)
		if err == nil {
			log.Printf("Successfully connected to Kafka")
			break
		}
		log.Printf("Failed to connect to Kafka (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create producer after %d attempts: %v", maxRetries, err)
	}

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		producer.Close()
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
	// Create topic if it doesn't exist
	admin, err := sarama.NewClusterAdmin(s.brokers, nil)
	if err != nil {
		log.Printf("Warning: Could not create Kafka admin client: %v", err)
	} else {
		log.Printf("Creating Kafka topic: %s", clickTopic)
		err = admin.CreateTopic(clickTopic, &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		if err != nil {
			if err == sarama.ErrTopicAlreadyExists {
				log.Printf("Topic %s already exists", clickTopic)
			} else {
				log.Printf("Warning: Could not create topic: %v", err)
			}
		} else {
			log.Printf("Successfully created topic: %s", clickTopic)
		}
		admin.Close()
	}

	// Get topic partitions
	partitions, err := s.consumer.Partitions(clickTopic)
	if err != nil {
		return fmt.Errorf("failed to get partitions: %v", err)
	}

	if len(partitions) == 0 {
		return fmt.Errorf("no partitions found for topic: %s", clickTopic)
	}

	log.Printf("Starting consumer for topic: %s, partition: %d", clickTopic, partitions[0])
	partitionConsumer, err := s.consumer.ConsumePartition(clickTopic, partitions[0], sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("failed to create partition consumer: %v", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		log.Printf("Consumer started, waiting for messages...")
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
				log.Printf("Consumer shutting down...")
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
