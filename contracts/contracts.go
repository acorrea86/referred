package contracts

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

// NeoSession is the interface for the neo4j session
type NeoSession interface {
	ReadTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error)
	WriteTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error)
}

// AppConsumer is the interface for the kafka consumer
type AppConsumer interface {
	StartConsumer()
}

// AppProducer is the interface for the kafka producer
type AppProducer interface {
	SendMessage(message []byte, topic string) error
}
