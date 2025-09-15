package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

// NewProducer creates a new Kafka producer with the given brokers and topic.
// It uses the LeastBytes balancer to distribute messages.
// The caller is responsible for closing the producer when it's no longer needed.
// Note: This is a simplified example. In a production scenario, you might want to
// configure additional settings like batch size, timeouts, retries, etc.
func NewProducer(brokers []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,                   // nombre de messages max par batch
		BatchTimeout: 10 * time.Millisecond, // délai max avant envoi du batch
		RequiredAcks: kafka.RequireAll,      // attendre ack de tous les brokers (sécurité)
		Async:        false,                 // true = meilleur débit, false = meilleure garantie
		Compression:  kafka.Snappy,          // compression efficace et rapide
		// Retry & Timeout settings
		// Retries: par défaut géré par le Writer en fonction des erreurs
		// mais on peut limiter avec WriteTimeout et ReadTimeout
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		// Logger: on peut brancher un logger custom ici
	}
}

// NewConsumer creates a new Kafka consumer with the given brokers, topic, and group ID.
// The caller is responsible for closing the consumer when it's no longer needed.
// Note: This is a simplified example. In a production scenario, you might want to
// configure additional settings like offsets, timeouts, etc.
func NewConsumer(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,             // 10KB
		MaxBytes:       10e6,             // 10MB
		CommitInterval: time.Second,      // commit auto chaque seconde
		MaxWait:        500 * time.Millisecond, // délai max d’attente des messages
		ReadLagInterval: -1,              // désactiver le suivi de lag si pas nécessaire
	})
}

// CloseProducer closes the given Kafka producer.
func CloseProducer(producer *kafka.Writer) error {
	return producer.Close()
}

// CloseConsumer closes the given Kafka consumer.
func CloseConsumer(consumer *kafka.Reader) error {
	return consumer.Close()
}
