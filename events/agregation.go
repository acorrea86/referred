package events

import "github.com/confluentinc/confluent-kafka-go/kafka"

type KafkaProducer struct {
	Producer *kafka.Producer
}

func (k *KafkaProducer) SendMessage(message []byte, topic string) error {
	return k.Producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          message,
		},
		make(chan kafka.Event, 10000),
	)
}

func NewProducer(producer *kafka.Producer) *KafkaProducer {
	return &KafkaProducer{Producer: producer}
}
