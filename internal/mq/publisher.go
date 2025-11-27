package mq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher keeps the rabbitmq connection, channel and exchange name.
type Publisher struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	exchange string
}

// NewPublisher conects to rabbitmq, opens a channel and declares an exchange.
func NewPublisher(url, exchange string) (*Publisher, error) {
	// connect to rabbitmq server
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	// opena channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}
	// declare a fanout exchange
	err = ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	log.Println("Publisher connected, exchange:", exchange)

	return &Publisher{
		conn:     conn,
		ch:       ch,
		exchange: exchange,
	}, nil
}

// PublishPangram sends a simple text message saying which pangram was found.
func (p *Publisher) PublishPangram(word string) error {
	if p == nil || p.ch == nil {
		return fmt.Errorf("publisher not initialized")
	}

	body := fmt.Sprintf("pangram found: %s", word)

	// use a context with timeout so publish doesn't hang forever
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.ch.PublishWithContext(ctx, p.exchange, "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	},
	)
	if err != nil {
		return fmt.Errorf("publish pangram: %w", err)
	}
	log.Println("RabbitMq published event:", body)
	return nil
}

// Close channel and connection
func (p *Publisher) Close() {
	if p == nil {
		return
	}
	if p.ch != nil {
		_ = p.ch.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}
