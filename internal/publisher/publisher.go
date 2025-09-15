package publisher

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cedricleblond35/HTTPserver/internal/kafka" // ton package interne
	k "github.com/segmentio/kafka-go"                      // alias pour le package segmentio
)

// Processor lit depuis un channel de Topic et publie chaque message dans Kafka.
func Processor(ctx context.Context, writer *k.Writer, input <-chan kafka.TopicMessage) error {
	errCh := make(chan error) // channel pour recevoir les erreurs d'envoi

	// Boucle principale
	for {
		select {
		case <-ctx.Done():
			return ctx.Err() // arrêt propre si le contexte est annulé
		case t := <-input:
			ProduceTopicAsync(ctx, writer, t, errCh)
		case err := <-errCh:
			// Gérer les erreurs d'envoi ici (log, métriques, etc.)
			return err
		}
	}
}

// ProduceTopicAsync sérialise un Topic et l’envoie de façon asynchrone à Kafka.
// errCh est un channel fourni par le caller pour recevoir les erreurs.
// Si errCh == nil, les erreurs sont ignorées.
func ProduceTopicAsync(ctx context.Context, writer *k.Writer, topic kafka.TopicMessage, errCh chan<- error) {
	// Sérialisation JSON
	data, err := json.Marshal(topic)
	if err != nil {
		if errCh != nil {
			errCh <- err
		}
		return
	}

	// Envoi asynchrone (goroutine)
	go func() {
		c, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = writer.WriteMessages(c, k.Message{
			Key:   []byte(topic.Name), // clé = partitionnement par nom
			Value: data,
			Time:  time.Now(),
		})

		if err != nil && errCh != nil {
			errCh <- err // on informe le caller
		}
	}()

}
