package kafkas

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"nav-receive-go/configs"
	"nav-receive-go/domains"
	"nav-receive-go/global"
	"nav-receive-go/services"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	defaultKafkaGroupID = "nav-receive-go"
	kafkaDefaultPort    = ":9092"
	kafkaReadTimeout    = 10 * time.Second
	kafkaCommitTimeout  = 5 * time.Second
	kafkaRetryDelay     = 5 * time.Second
)

var rtloggingConsumerOnce sync.Once

func StartRtloggingConsumer() {
	rtloggingConsumerOnce.Do(func() {
		kafkaCfg := global.NAV_CONFIG.Kafka
		if len(kafkaCfg.Servers) == 0 || strings.TrimSpace(kafkaCfg.RtloggingTopic) == "" {
			log.Println("kafka rtlogging consumer disabled: servers or rtlogging-topic is empty")
			return
		}
		go consumeRtloggingData(kafkaCfg)
	})
}

func consumeRtloggingData(kafkaCfg configs.Kafka) {
	groupID := strings.TrimSpace(kafkaCfg.GroupID)
	if groupID == "" {
		groupID = defaultKafkaGroupID
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        normalizeKafkaServers(kafkaCfg.Servers),
		Topic:          kafkaCfg.RtloggingTopic,
		GroupID:        groupID,
		MinBytes:       1,
		MaxBytes:       10 * 1024 * 1024,
		MaxWait:        time.Second,
		CommitInterval: 0,
	})
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("kafka rtlogging consumer close error: %v", err)
		}
	}()

	log.Printf("✅kafka rtlogging consumer started: topic=%s group=%s brokers=%v", kafkaCfg.RtloggingTopic, groupID, normalizeKafkaServers(kafkaCfg.Servers))
	for {
		ctx, cancel := context.WithTimeout(context.Background(), kafkaReadTimeout)
		msg, err := reader.FetchMessage(ctx)
		cancel()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				continue
			}
			log.Printf("❌kafka rtlogging fetch message error: %v", err)
			time.Sleep(kafkaRetryDelay)
			continue
		}

		if err = handleRtloggingMessage(msg); err != nil {
			log.Printf("✅kafka rtlogging handle message error: key=%s offset=%d err=%v", string(msg.Key), msg.Offset, err)
		}
		commitCtx, commitCancel := context.WithTimeout(context.Background(), kafkaCommitTimeout)
		if err = reader.CommitMessages(commitCtx, msg); err != nil {
			log.Printf("❌kafka rtlogging commit message error: key=%s offset=%d err=%v", string(msg.Key), msg.Offset, err)
		}
		commitCancel()
	}
}

func handleRtloggingMessage(msg kafka.Message) error {
	var rtloggingData domains.RtloggingData[domains.GP]
	if err := json.Unmarshal(msg.Value, &rtloggingData); err != nil {
		return fmt.Errorf("❌unmarshal rtlogging data: %w", err)
	}

	gpsInitial := strings.TrimSpace(rtloggingData.Data.GpsInitial)
	if gpsInitial == "" {
		return nil
	}
	rtcmBytes, err := base64.StdEncoding.DecodeString(gpsInitial)
	if err != nil {
		return fmt.Errorf("❌decode gpsInitial: %w", err)
	}

	sncode := resolveRtloggingSncode(rtloggingData, msg.Key)
	if sncode == "" {
		return fmt.Errorf("❌empty sncode")
	}
	services.WriteRtcmFile(sncode, parseRtloggingTime(rtloggingData.Time), rtcmBytes)
	return nil
}

func resolveRtloggingSncode(rtloggingData domains.RtloggingData[domains.GP], key []byte) string {
	candidates := []string{
		rtloggingData.Device,
		rtloggingData.Data.Sncode,
		rtloggingData.DeviceId,
		string(key),
	}
	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}
		if services.IsDeviceSncodeCached(candidate) {
			return candidate
		}
	}
	for _, candidate := range candidates {
		if candidate = strings.TrimSpace(candidate); candidate != "" {
			return candidate
		}
	}
	return ""
}

func parseRtloggingTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Now().UTC()
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05"} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC()
		}
	}
	return time.Now().UTC()
}

func normalizeKafkaServers(servers []string) []string {
	result := make([]string, 0, len(servers))
	for _, server := range servers {
		server = strings.TrimSpace(server)
		if server == "" {
			continue
		}
		if !strings.Contains(server, ":") {
			server += kafkaDefaultPort
		}
		result = append(result, server)
	}
	return result
}
