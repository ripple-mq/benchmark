package tape

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/ripple-mq/benchmark/tape/kafka"
	"github.com/ripple-mq/benchmark/tape/ripple"
)

func RunRipple(msgCount int, messageSizeKB int, batchSize int) {
	rippleRes := ripple.BenchmarkConsumerThroughput(msgCount, messageSizeKB, batchSize)
	AppendMarkdown("./benchmark.md", rippleRes)
}

func RunKafka(msgCount int, messageSizeKB int, batchSize int) {
	kafkaRes := kafka.BenchmarkConsumerThroughput(msgCount, messageSizeKB, batchSize)
	AppendMarkdown("./benchmark.md", kafkaRes)
}

func PrettyPrint(p any) {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(data))
}

func AppendMarkdown(filePath string, data any) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("AppendMarkdown only supports struct types")
	}

	t := v.Type()

	maxFieldLen := len("Field")
	maxValueLen := len("Value")
	fields := make([][2]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get("json")
		if fieldName == "" {
			fieldName = field.Name
		}
		valStr := fmt.Sprintf("%v", v.Field(i).Interface())
		fields[i] = [2]string{fieldName, valStr}

		if len(fieldName) > maxFieldLen {
			maxFieldLen = len(fieldName)
		}
		if len(valStr) > maxValueLen {
			maxValueLen = len(valStr)
		}
	}

	host, _ := os.Hostname()
	currentUser, _ := user.Current()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	sysInfo := [][2]string{
		{"OS", runtime.GOOS},
		{"Arch", runtime.GOARCH},
		{"Hostname", host},
		{"User", currentUser.Username},
		{"Go Version", runtime.Version()},
		{"CPU Cores", fmt.Sprintf("%d", runtime.NumCPU())},
		{"Total Allocated Memory (MB)", fmt.Sprintf("%.2f", float64(memStats.TotalAlloc)/1024/1024)},
		{"Heap Allocated Memory (MB)", fmt.Sprintf("%.2f", float64(memStats.HeapAlloc)/1024/1024)},
	}

	for _, item := range sysInfo {
		if len(item[0]) > maxFieldLen {
			maxFieldLen = len(item[0])
		}
		if len(item[1]) > maxValueLen {
			maxValueLen = len(item[1])
		}
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("## 📊 Report - %s\n\n", time.Now().Format(time.RFC1123)))

	headerFmt := fmt.Sprintf("| %%-%ds | %%-%ds |\n", maxFieldLen, maxValueLen)
	sepFmt := fmt.Sprintf("|%%-%ds-|%%-%ds-|\n", maxFieldLen+1, maxValueLen+1)
	rowFmt := headerFmt

	b.WriteString("### 🛠 System Metadata\n\n")
	b.WriteString(fmt.Sprintf(headerFmt, "Field", "Value"))
	b.WriteString(fmt.Sprintf(sepFmt, strings.Repeat("-", maxFieldLen), strings.Repeat("-", maxValueLen)))
	for _, item := range sysInfo {
		b.WriteString(fmt.Sprintf(rowFmt, item[0], item[1]))
	}
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf(sepFmt, strings.Repeat("-", maxFieldLen), strings.Repeat("-", maxValueLen)))
	for _, field := range fields {
		b.WriteString(fmt.Sprintf(rowFmt, field[0], field[1]))
	}
	b.WriteString("\n---\n\n")

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open markdown file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(b.String()); err != nil {
		return fmt.Errorf("failed to write markdown: %w", err)
	}

	return nil
}
