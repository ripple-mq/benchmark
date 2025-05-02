// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"strings"
// 	"time"

// 	"github.com/segmentio/kafka-go"
// )

// const (
// 	topic      = "test-topic"
// 	brokerAddr = "localhost:9092"
// 	messageCnt = 10000
// )

// func createTopic() {
// 	conn, err := kafka.Dial("tcp", brokerAddr)
// 	if err != nil {
// 		log.Fatal("Failed to dial Kafka:", err)
// 	}
// 	defer conn.Close()

// 	controller, err := conn.Controller()
// 	if err != nil {
// 		log.Fatal("Failed to get controller:", err)
// 	}
// 	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
// 	if err != nil {
// 		log.Fatal("Failed to connect to controller:", err)
// 	}
// 	defer controllerConn.Close()

// 	topicConfigs := []kafka.TopicConfig{
// 		{
// 			Topic:             topic,
// 			NumPartitions:     1,
// 			ReplicationFactor: 1,
// 		},
// 	}

// 	err = controllerConn.CreateTopics(topicConfigs...)
// 	if err != nil {
// 		log.Fatal("Failed to create topic:", err)
// 	}

// 	fmt.Println("✅ Topic created (if not already existing)")
// }

// func produce() {
// 	w := kafka.Writer{
// 		Addr:         kafka.TCP(brokerAddr),
// 		Topic:        topic,
// 		Balancer:     &kafka.LeastBytes{},
// 		BatchSize:    100, // You can tweak this
// 		BatchTimeout: 10 * time.Millisecond,
// 		Async:        false,
// 	}

// 	defer w.Close()

// 	fmt.Println("🚀 Producing messages...")

// 	start := time.Now()
// 	totalBytes := 0

// 	messages := make([]kafka.Message, 0, 100)

// 	for i := 0; i < messageCnt; i++ {
// 		msg := kafka.Message{
// 			Key:   []byte(fmt.Sprintf("Key-%d", i)),
// 			Value: []byte(strings.Repeat("x", 1024)),
// 		}

// 		messages = append(messages, msg)
// 		totalBytes += len(msg.Key) + len(msg.Value)

// 		// Send in batches of 100 messages
// 		if len(messages) == 100 {
// 			if err := w.WriteMessages(context.Background(), messages...); err != nil {
// 				log.Fatalf("Failed to write messages: %v", err)
// 			}
// 			messages = messages[:0] // Reset slice
// 		}
// 	}

// 	// Flush any remaining messages
// 	if len(messages) > 0 {
// 		if err := w.WriteMessages(context.Background(), messages...); err != nil {
// 			log.Fatalf("Failed to write remaining messages: %v", err)
// 		}
// 	}

// 	elapsed := time.Since(start)
// 	mb := float64(totalBytes) / (1024 * 1024)
// 	throughput := mb / elapsed.Seconds()

// 	fmt.Printf("✅ Finished producing 10K messages in %s (%.2f MB, %.2f MB/s)\n", elapsed, mb, throughput)
// }

// func consume() {
// 	r := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers:         []string{brokerAddr},
// 		Topic:           topic,
// 		GroupID:         "test-group",
// 		MinBytes:        10e3,  // 10KB
// 		MaxBytes:        100e6, // 10MB
// 		MaxWait:         50 * time.Millisecond,
// 		ReadLagInterval: -1,
// 		QueueCapacity:   1000,
// 	})

// 	defer r.Close()

// 	fmt.Println("🧲 Consuming messages...")

// 	start := time.Now()
// 	count := 0
// 	totalBytes := 0

// 	for count < messageCnt {
// 		m, err := r.ReadMessage(context.Background())
// 		if err != nil {
// 			log.Fatalf("Failed to read message: %v", err)
// 		}
// 		count++
// 		totalBytes += len(m.Key) + len(m.Value)

// 		if count%1000 == 0 {
// 			fmt.Printf("Consumed %d messages...\n", count)
// 		}
// 	}

// 	elapsed := time.Since(start)
// 	mb := float64(totalBytes) / (1024 * 1024)
// 	throughput := mb / elapsed.Seconds()

// 	fmt.Printf("✅ Consumed 10K messages in %s (%.2f MB, %.2f MB/s)\n", elapsed, mb, throughput)
// }

//	func main() {
//		createTopic()
//		produce()
//		consume()
//	}

package main

import (
	"time"

	"github.com/briandowns/spinner"
)

func main() {
	s := spinner.New(spinner.CharSets[11], 120*time.Millisecond)
	s.Suffix = " Waiting for event..."
	s.FinalMSG = "✅ Event completed successfully!\n"
	s.Start()

	// Simulate waiting for an event
	waitForEvent()

	s.Stop()
}

func waitForEvent() {
	time.Sleep(3 * time.Second)
}
