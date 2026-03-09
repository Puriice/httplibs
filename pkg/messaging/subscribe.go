package messaging

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitListener struct {
	queueName string
	channel   *amqp.Channel
}

func (r *RabbitMQ) NewListener(queueName string, keys ...string) (*RabbitListener, error) {
	q, err := r.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		err := r.Channel.QueueBind(
			q.Name,
			key,
			r.Exchange,
			false,
			nil,
		)

		if err != nil {
			return nil, err
		}
	}

	return &RabbitListener{
		queueName: q.Name,
		channel:   r.Channel,
	}, nil
}

func (l *RabbitListener) Subscribe(handler func([]byte) error) error {
	err := l.channel.Qos(10, 0, false)
	if err != nil {
		return err
	}

	msgs, err := l.channel.Consume(
		l.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	for msg := range msgs {
		if err := handler(msg.Body); err != nil {
			log.Println(err)

			msg.Nack(false, true)

			continue
		}

		msg.Ack(true)
	}

	return nil
}
