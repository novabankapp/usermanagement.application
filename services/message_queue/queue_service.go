package message_queue

import (
	"context"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	kafka_go "github.com/segmentio/kafka-go"
	"time"
)

type MessageQueue interface {
	PublishMessage(ctx context.Context, message []byte, key string, topic string) (bool, error)
}

type KafkaMessageQueue struct {
	producer kafkaClient.Producer
}

func NewKafkaMessageQueue(producer kafkaClient.Producer) MessageQueue {
	return &KafkaMessageQueue{
		producer: producer,
	}
}
func (k KafkaMessageQueue) PublishMessage(ctx context.Context, message []byte, key string, topic string) (bool, error) {
	mes := kafka_go.Message{
		Topic: topic, //r.cfg.KafkaTopics.UserCreated.TopicName
		Value: message,
		Time:  time.Now().UTC(),
		Key:   []byte(key),
	}
	if error := k.producer.PublishMessage(ctx, mes); error != nil {
		return false, error
	}
	return true, nil

}
