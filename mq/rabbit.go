package mq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SendToRabbit(frame []byte, cameraID string) {
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		log.Println("RabbitMQ conn error:", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("RabbitMQ channel error:", err)
		return
	}
	defer ch.Close()

	queue := "cctv_frames"
	_, _ = ch.QueueDeclare(queue, true, false, false, false, nil)

	_ = ch.Publish("", queue, false, false, amqp.Publishing{
		ContentType: "image/jpeg",
		Body:        frame,
	})
}
