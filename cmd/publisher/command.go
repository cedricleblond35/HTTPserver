package publisher

import (
	"context"
	"log"
	"time"

	"github.com/cedricleblond35/HTTPserver/internal/kafka"
	"github.com/cedricleblond35/HTTPserver/internal/publisher"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "collector",
	Short: "This is the todo collector",
	Long:  "This binary will poll kafka todo topic, and collect them",
	RunE:  action,
}

func action(cmd *cobra.Command, _ []string) (err error) {

	ctx := context.Background()

	// // Création du producer Kafka
	writer := kafka.NewProducer([]string{"localhost:9092"}, "my-topic")
	defer kafka.CloseProducer(writer)

	// Channel pour envoyer des Topics à Processor
	input := make(chan kafka.TopicMessage)

	// Lancement du Processor dans une goroutine
	go func() {
		if err := publisher.Processor(ctx, writer, input); err != nil {
			log.Printf("error producing topic: %v", err)
		}
	}()

	// Channel pour récupérer les erreurs de production
	errCh := make(chan error, 10) // buffer pour ne pas bloquer la goroutine

	// Exemple d’envoi de message via ProduceTopicAsync direct
	topic1 := kafka.TopicMessage{
		Name:    "direct",
		Message: "message direct avec ProduceTopicAsync",
		Date:    time.Now(),
	}
	publisher.ProduceTopicAsync(ctx, writer, topic1, errCh)
	// Lecture de l’erreur si elle arrive
	select {
	case err := <-errCh:
		log.Printf("error producing topic: %v", err)
	case <-time.After(1 * time.Second):
		// pas d'erreur dans la seconde
	}	

	// Exemple d’envoi de message via Processor (channel)
	topic2 := kafka.TopicMessage{
		Name:    "processor",
		Message: "message envoyé via Processor",
		Date:    time.Now(),
	}
	input <- topic2

	// Attendre un peu pour laisser les goroutines finir
	time.Sleep(2 * time.Second)

	return
}
