package messaging

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitListener struct {
	queue   *amqp.Queue
	channel *amqp.Channel
}

func (r *RabbitMQ) NewListener(queueName string, keys ...string) (*RabbitListener, error) {
	q, err := r.Channel.QueueDeclare(
		queueName,
		false,
		false,
		true,
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
		queue:   &q,
		channel: r.Channel,
	}, nil
}

func (l *RabbitListener) Subscribe(handler func([]byte) error) error {
	msgs, err := l.channel.Consume(
		l.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		for data := range msgs {
			err = handler(data.Body)

			if err != nil {
				log.Println(err)
				continue
			}
		}
	}()

	return nil
}
