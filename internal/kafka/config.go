package kafka

import "time"

type TopicMessage struct {
	Name    string    `json:"name"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}
