package events

import (
	"blumer-ms-refers/model"
	"errors"
	"fmt"
	"reflect"

	"blumer-ms-refers/repository"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaReducer is the kafka Consumer
type KafkaReducer struct {
	Repository *repository.Repository
	Consumer   *kafka.Consumer
	Producer   *KafkaProducer
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

	err := k.Consumer.SubscribeTopics([]string{
		"ms-profile-create",
		"ms-profile-edit",
		"ms-profile-delete",
		"ms-wallet-validation",
	}, nil)
	if err != nil {
		fmt.Printf("error subscribing topics %v \n", err.Error())
	}

	for {
		msg, err := k.Consumer.ReadMessage(-1)
		if err == nil {
			profile, err := model.DecodeProfile(msg.Value)
			if err != nil {
				fmt.Printf("Can't decode data: %v (%v)\n", err, msg)
				break
			}
			switch *msg.TopicPartition.Topic {
			case "ms-profile-create":
				err = k.saveProfile(profile)
				if err != nil {
					fmt.Printf("Can't save profile: %v (%v)\n", err, profile)
				}
				break
			case "ms-profile-edit":
				err = k.updateProfile(profile)
				if err != nil {
					fmt.Printf("Can't update profile: %v (%v)\n", err, profile)
				}
				break
			case "ms-profile-delete":
				err = k.deleteProfile(profile)
				if err != nil {
					fmt.Printf("Can't delete profile: %v (%v)\n", err, profile)
				}
				break
			case "ms-wallet-validation":
				userID, err := model.DecodeUserID(msg.Value)
				if err != nil {
					fmt.Printf("Can't decode data: %v (%v)\n", err, msg)
					break
				}

				err = k.validateReferred(*userID)
				if err != nil {
					fmt.Printf("Can't validate referred: %v (%v)\n", err, userID)
				}
				break
			}

		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}

func (k *KafkaReducer) saveProfile(data *model.Profile) error {
	profile, err := k.Repository.Find(data.UserID)
	if err != nil {
		return err
	}
	if profile != nil {
		return nil
	}

	return k.Repository.Save(data)
}

func (k *KafkaReducer) updateProfile(data *model.Profile) error {
	profile, err := k.Repository.Find(data.UserID)
	if err != nil {
		return err
	}
	if profile == nil {
		return errors.New("profile not found")
	}

	if reflect.DeepEqual(profile, data) {
		return nil
	}

	return k.Repository.Update(data)

}

func (k *KafkaReducer) deleteProfile(data *model.Profile) error {
	profile, err := k.Repository.Find(data.UserID)
	if err != nil {
		return err
	}

	if profile == nil {
		return errors.New("profile not found")
	}

	return k.Repository.Delete(data.UserID)
}

func (k *KafkaReducer) validateReferred(userID string) error {
	profile, err := k.Repository.Find(userID)
	if err != nil {
		return err
	}

	if profile == nil {
		return errors.New("profile not found")
	}

	referrer, err := k.Repository.GetReferrer(userID)
	if err != nil {
		return err
	}

	msj, err := model.EncodeWalletReward(*referrer)
	if err != nil {
		return err
	}

	return k.Producer.SendMessage(msj, "ms-wallet-reward")

}

// NewKafkaReducer creates a new kafka Consumer
func NewKafkaReducer(
	repository *repository.Repository,
	consumer *kafka.Consumer,
	producer *KafkaProducer,
) *KafkaReducer {
	return &KafkaReducer{
		Repository: repository,
		Consumer:   consumer,
		Producer:   producer,
	}
}
