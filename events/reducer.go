package events

import (
	"fmt"

	"blumer-ms-refers/repository"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaReducer is the kafka Consumer
type KafkaReducer struct {
	Repository *repository.Repository
	Consumer   *kafka.Consumer
}

// StartConsumer starts the kafka Consumer
func (k *KafkaReducer) StartConsumer() {
	fmt.Print("starting Consumer")
	defer func(consumer *kafka.Consumer) {
		err := consumer.Close()
		if err != nil {
			fmt.Printf("error closing Consumer %v \n", err)
		}
	}(k.Consumer)

	err := k.Consumer.SubscribeTopics([]string{"ms-profile-create", "ms-profile-edit", "ms-profile-delete"}, nil)
	if err != nil {
		fmt.Printf("error subscribing topics %v \n", err.Error())
	}

	for {
		msg, err := k.Consumer.ReadMessage(-1)
		if err == nil {
			switch *msg.TopicPartition.Topic {
			case "ms-profile-create":
				break
			case "ms-profile-edit":
				break
			case "ms-profile-delete":
				break
			}

		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}

// NewKafkaReducer creates a new kafka Consumer
func NewKafkaReducer(
	repository *repository.Repository,
	consumer *kafka.Consumer,
) *KafkaReducer {
	return &KafkaReducer{
		Repository: repository,
		Consumer:   consumer,
	}
}
