package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ripple-mq/benchmark/pkg/utils/pen"
	"github.com/segmentio/kafka-go"
)

type ConsumeThroughput struct {
	Desc                 string  `json:"description"`
	Broker               string  `json:"broker"`
	ReadBatchSizeInBytes int     `json:"read_batch_size"`
	MessageSizeInBytes   int     `json:"message_size_in_bytes"`
	TotalDataInBytes     float64 `json:"total_data_in_bytes"`
	ValueInMBPerSec      float64 `json:"value_in_mb_per_sec"`
}

const (
	topic      = "test-topic"
	brokerAddr = "localhost:9092"
)

func createTopic() {
	conn, err := kafka.Dial("tcp", brokerAddr)
	if err != nil {
		log.Fatal("Failed to dial Kafka:", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatal("Failed to get controller:", err)
	}
	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		log.Fatal("Failed to connect to controller:", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Fatal("Failed to create topic:", err)
	}
}

func produce(messageCnt int, messageSizeKB int) {
	w := kafka.Writer{
		Addr:         kafka.TCP(brokerAddr),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1000,
		BatchTimeout: 50 * time.Millisecond,
		Async:        true,
		RequiredAcks: kafka.RequireNone,
		Compression:  kafka.Snappy,
	}

	defer w.Close()

	totalBytes := 0
	messages := make([]kafka.Message, 0, 100)

	l := pen.Loader("Producing 10K messages ")
	for i := 0; i < messageCnt; i++ {
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("Key-%d", i)),
			Value: []byte(strings.Repeat("x", 1024*messageSizeKB)),
		}

		messages = append(messages, msg)
		totalBytes += len(msg.Key) + len(msg.Value)

		if len(messages) == 100 {
			if err := w.WriteMessages(context.Background(), messages...); err != nil {
				log.Fatalf("Failed to write messages: %v", err)
			}
			messages = messages[:0]
		}
	}

	if len(messages) > 0 {
		if err := w.WriteMessages(context.Background(), messages...); err != nil {
			log.Fatalf("Failed to write remaining messages: %v", err)
		}
	}

	pen.Complete(l, "Produced ~1M messages successfully ")
}

func consume(messageCnt int, messageSizeKB int, batchSize int) ConsumeThroughput {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         []string{brokerAddr},
		Topic:           topic,
		GroupID:         "test-group",
		MinBytes:        10e3,
		MaxBytes:        batchSize,
		MaxWait:         50 * time.Millisecond,
		ReadLagInterval: -1,
		QueueCapacity:   1000,
	})

	defer r.Close()

	start := time.Now()
	count := 0
	totalBytes := 0

	for count < messageCnt {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			log.Fatalf("Failed to fetch message: %v", err)
		}

		totalBytes += len(m.Key) + len(m.Value)
		count++

		if count%100 == 0 {
			if err := r.CommitMessages(context.Background(), m); err != nil {
				log.Printf("commit failed: %v", err)
			}
		}
	}

	elapsed := time.Since(start)
	mb := float64(totalBytes) / (1024 * 1024)
	throughput := mb / elapsed.Seconds()

	return ConsumeThroughput{
		Desc:                 "network I/O througput",
		Broker:               "kafka",
		ReadBatchSizeInBytes: r.Config().MaxBytes,
		MessageSizeInBytes:   1024 * messageSizeKB,
		TotalDataInBytes:     float64(messageCnt * 1024.0),
		ValueInMBPerSec:      throughput,
	}
}

func BenchmarkConsumerThroughput(messageCnt int, messageSizeKB int, batchSize int) ConsumeThroughput {
	createTopic()
	produce(messageCnt, messageSizeKB)
	return consume(messageCnt, messageSizeKB, batchSize)
}
