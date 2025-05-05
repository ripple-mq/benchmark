package prometheus

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ripple-mq/go-client/wave"
)

var (
	cnt = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "consumer_network_i/o_througput",
			Help: "Consumer network I/O throughput",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(cnt)
}

func ConsumerThroughput() {
	client := wave.NewClient[string](wave.Config{
		Topic:         "topic-X",
		Bucket:        "topic-Y",
		Brokers:       []string{":8080"},
		ReadBatchSize: 4,
	})

	client.Create()

	prodCh := client.Writer()

	log.Info("Producing nearly a million messages of 1kb")
	go func() {
		for i := 0; i < 1000100; i++ {
			prodCh <- strings.Repeat("x", 1024)
		}
	}()

	time.Sleep(1 * time.Minute)
	fmt.Println("Producing Done")

	go func() {
		for {
			consCh := client.Reader()
			const messageCount = 10000
			const messageSizeKB = 1
			start := time.Now()

			for i := 0; i < messageCount; i++ {
				<-consCh
			}

			elapsed := time.Since(start).Seconds()
			totalMB := float64(messageCount*messageSizeKB) / 1024.0
			throughput := totalMB / elapsed

			fmt.Printf("Consumed %d messages of %d KB each in %.3f seconds\n", messageCount, messageSizeKB, elapsed)
			fmt.Printf("Consumer throughput: %.2f MB/s\n", throughput)

			cnt.WithLabelValues("network IO").Set(throughput)
		}
	}()

}
