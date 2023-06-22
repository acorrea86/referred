package di

import (
	"blumer-ms-refers/contracts"
	"blumer-ms-refers/handler"
	"blumer-ms-refers/repository"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"os"

	"blumer-ms-refers/events"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func providerRepository(session contracts.NeoSession) (*repository.Repository, error) {
	return repository.NewRepository(session), nil
}

func providerNeoSession() (contracts.NeoSession, error) {
	neo4jUri := os.Getenv("NEO4J_URI")
	if neo4jUri == "" {
		return nil, errors.New("env var NEO4J_URI is not defined")
	}

	neo4jdb := os.Getenv("NEO4J_DATABASE")
	if neo4jdb == "" {
		return nil, errors.New("env var NEO4J_DATABASE is not defined")
	}

	neo4jUsername := os.Getenv("NEO4J_USERNAME")
	if neo4jUsername == "" {
		return nil, errors.New("env var NEO4J_USERNAME is not defined")
	}

	neo4jPassword := os.Getenv("NEO4J_PASSWORD")
	if neo4jPassword == "" {
		return nil, errors.New("env var NEO4J_PASSWORD is not defined")
	}

	driver, err := neo4j.NewDriver(neo4jUri, neo4j.BasicAuth(neo4jUsername, neo4jPassword, ""))
	if err != nil {
		return nil, err
	}
	session := driver.NewSession(neo4j.SessionConfig{
		DatabaseName: neo4jdb,
	})

	return session, nil
}

func providerKafkaConsumer() (*kafka.Consumer, error) {
	groupId := os.Getenv("KAFKA_GROUP_ID")
	if groupId == "" {
		return nil, errors.New("env var KAFKA_GROUP_ID is not defined")
	}

	bootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	if bootstrapServers == "" {
		return nil, errors.New("env var KAFKA_BOOTSTRAP_SERVER is not defined")
	}

	securityProtocol := os.Getenv("SECURITY_PROTOCOL")
	if securityProtocol == "" {
		return nil, errors.New("env var SECURITY_PROTOCOL is not defined")
	}

	saslMechanisms := os.Getenv("SASL_MECHANISMS")
	if saslMechanisms == "" {
		return nil, errors.New("env var SASL_MECHANISMS is not defined")
	}

	saslUsername := os.Getenv("SASL_USERNAME")
	if saslUsername == "" {
		return nil, errors.New("env var SASL_USERNAME is not defined")
	}

	saslPassword := os.Getenv("SASL_PASSWORD")
	if saslPassword == "" {
		return nil, errors.New("env var SASL_PASSWORD is not defined")
	}

	data, err := kafka.NewConsumer(&kafka.ConfigMap{
		"group.id":          groupId,
		"bootstrap.servers": bootstrapServers,
		"security.protocol": securityProtocol,
		"sasl.mechanisms":   saslMechanisms,
		"sasl.username":     saslUsername,
		"sasl.password":     saslPassword,
		"auto.offset.reset": "earliest",
	})

	return data, err
}

func providerKafkaProducer() (*kafka.Producer, error) {
	bootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	if bootstrapServers == "" {
		return nil, errors.New("env var KAFKA_BOOTSTRAP_SERVER is not defined")
	}

	securityProtocol := os.Getenv("SECURITY_PROTOCOL")
	if securityProtocol == "" {
		return nil, errors.New("env var SECURITY_PROTOCOL is not defined")
	}

	saslMechanisms := os.Getenv("SASL_MECHANISMS")
	if saslMechanisms == "" {
		return nil, errors.New("env var SASL_MECHANISMS is not defined")
	}

	saslUsername := os.Getenv("SASL_USERNAME")
	if saslUsername == "" {
		return nil, errors.New("env var SASL_USERNAME is not defined")
	}

	saslPassword := os.Getenv("SASL_PASSWORD")
	if saslPassword == "" {
		return nil, errors.New("env var SASL_PASSWORD is not defined")
	}

	messageTimeoutMS := os.Getenv("MESSAGE_TIMEOUT_MS")
	if messageTimeoutMS == "" {
		return nil, errors.New("env var MESSAGE_TIMEOUT_MS is not defined")
	}
	data, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  bootstrapServers,
		"security.protocol":  securityProtocol,
		"sasl.mechanisms":    saslMechanisms,
		"sasl.username":      saslUsername,
		"sasl.password":      saslPassword,
		"message.timeout.ms": messageTimeoutMS,
	})

	return data, err
}

func providerReducer(
	repository *repository.Repository,
	consumer *kafka.Consumer,
	producer *events.KafkaProducer,
) *events.KafkaReducer {
	return events.NewKafkaReducer(repository, consumer, producer)
}

func providerAggregation(producer *kafka.Producer) *events.KafkaProducer {
	return events.NewProducer(producer)
}

func providerHandler(
	repository *repository.Repository,
	reducer *events.KafkaReducer,
	producer *events.KafkaProducer,
) (*handler.Handler, error) {
	graphqlPort := os.Getenv("GRAPHQL_PORT")
	if graphqlPort == "" {
		return nil, errors.New("env var GRAPHQL_PORT is not defined")
	}
	return handler.NewHandler(repository, reducer, producer, graphqlPort), nil
}
