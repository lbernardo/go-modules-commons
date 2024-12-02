package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type IRabbit interface {
	IsClosed() bool
	Refresh()
	PublishMessage(ctx context.Context, queueName string, messageId string, message any) error
	CloseConsumer(consumerTag string) error
	PublishExchange(exchangeName string, message any) error
	Subscribe(exchangeName string) (<-chan amqp.Delivery, error)
	CreateConsumer(name string, consumerTag string) (<-chan amqp.Delivery, error)
}

type RabbitMQ struct {
	channel   *amqp.Channel
	conn      *amqp.Connection
	queues    map[string]amqp.Queue
	exchanges map[string]bool
	logger    *zap.Logger
	cfg       *viper.Viper
}

func NewRabbitMQ(cfg *viper.Viper, logger *zap.Logger) IRabbit {
	r := &RabbitMQ{
		cfg:    cfg,
		logger: logger.Named("rabbitmq"),
	}
	if err := r.connect(); err != nil {
		r.logger.Fatal("failed to connect rabbitmq", zap.Error(err))
		return nil
	}
	return r
}

func (r *RabbitMQ) IsClosed() bool {
	if r.channel == nil {
		return true
	}
	return r.channel.IsClosed()
}

func (r *RabbitMQ) connect() error {
	conn, err := amqp.Dial(r.cfg.GetString("app.queue.host"))
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel in RabbitMQ: %v", err)
	}
	r.conn = conn
	r.channel = channel
	return nil
}

func (r *RabbitMQ) Refresh() {
	if r.IsClosed() {
		r.logger.Info("RabbitMQ connection is closed! Needs refresh...")
		if err := r.connect(); err != nil {
			r.logger.Fatal("failed to refresh RabbitMQ", zap.Error(err))
			return
		}
		r.logger.Info("RabbitMQ connection refreshed")
		return
	}
	r.logger.Info("RabbitMQ connection it's ok!")
}

func (r *RabbitMQ) registerQueue(name string) (amqp.Queue, error) {
	if name == "" {
		return amqp.Queue{}, errors.New("queue name is required")
	}
	queue, ok := r.queues[name]
	if !ok {
		q, err := r.channel.QueueDeclare(name,
			false,
			false,
			false,
			false,
			nil)
		return q, err
	}
	return queue, nil
}

func (r *RabbitMQ) registerExchange(name string) error {
	if name == "" {
		return errors.New("exchange name is required")
	}
	if _, ok := r.exchanges[name]; ok {
		return nil
	}
	err := r.channel.ExchangeDeclare(name, amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil)
	return err
}

func (r *RabbitMQ) PublishExchange(exchangeName string, message any) error {
	if err := r.registerExchange(exchangeName); err != nil {
		return err
	}
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}
	return r.channel.Publish(exchangeName, "",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (r *RabbitMQ) Subscribe(exchangeName string) (<-chan amqp.Delivery, error) {
	if err := r.registerExchange(exchangeName); err != nil {
		return nil, err
	}
	q, err := r.channel.QueueDeclare("",
		false,
		false,
		true,
		false,
		nil)
	if err != nil {
		return nil, fmt.Errorf("error declaring queue: %w", err)
	}
	if err := r.channel.QueueBind(q.Name, "", exchangeName, false, nil); err != nil {
		return nil, fmt.Errorf("error binding queue: %w", err)
	}
	return r.channel.Consume(q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, queueName string, messageId string, message any) error {
	queue, err := r.registerQueue(queueName)
	if err != nil {
		return err
	}
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}

	return r.channel.PublishWithContext(ctx, "", queue.Name,
		false,
		false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			MessageId:   messageId,
		})
}

func (r *RabbitMQ) CreateConsumer(name string, consumerTag string) (<-chan amqp.Delivery, error) {
	q, err := r.registerQueue(name)
	if err != nil {
		return nil, err
	}

	return r.channel.Consume(q.Name,
		consumerTag,
		true,
		false,
		false,
		false,
		nil)
}

func (r *RabbitMQ) CloseConsumer(consumerTag string) error {
	return r.channel.Cancel(consumerTag, false)
}
