package mq

import (
	"encoding/json"
	"fmt"
	"time"

	"ocean-marketing/internal/config"
	"ocean-marketing/internal/pkg/logger"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// Client 消息队列客户端
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     config.MQConfig
}

// Message 消息结构
type Message struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
	Retry     int                    `json:"retry"`
}

// NewClient 创建消息队列客户端
func NewClient(cfg config.MQConfig) (*Client, error) {
	client := &Client{cfg: cfg}

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

// connect 连接到RabbitMQ
func (c *Client) connect() error {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		c.cfg.Username, c.cfg.Password, c.cfg.Host, c.cfg.Port, c.cfg.Vhost)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		logger.Error("连接RabbitMQ失败", zap.Error(err))
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		logger.Error("创建RabbitMQ channel失败", zap.Error(err))
		conn.Close()
		return err
	}

	c.conn = conn
	c.channel = channel

	logger.Info("RabbitMQ连接成功")
	return nil
}

// Publish 发布消息
func (c *Client) Publish(exchange, routingKey string, message Message) error {
	message.Timestamp = time.Now().Unix()

	body, err := json.Marshal(message)
	if err != nil {
		logger.Error("序列化消息失败", zap.Error(err))
		return err
	}

	err = c.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		logger.Error("发布消息失败", zap.Error(err))
		return err
	}

	logger.Info("消息发布成功",
		zap.String("exchange", exchange),
		zap.String("routing_key", routingKey),
		zap.String("message_id", message.ID))

	return nil
}

// Subscribe 订阅消息
func (c *Client) Subscribe(queueName string, handler func(Message) error) error {
	// 声明队列
	_, err := c.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		logger.Error("声明队列失败", zap.Error(err))
		return err
	}

	// 设置QoS
	err = c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		logger.Error("设置QoS失败", zap.Error(err))
		return err
	}

	// 消费消息
	msgs, err := c.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		logger.Error("消费消息失败", zap.Error(err))
		return err
	}

	go func() {
		for d := range msgs {
			var message Message
			if err := json.Unmarshal(d.Body, &message); err != nil {
				logger.Error("反序列化消息失败", zap.Error(err))
				d.Nack(false, false)
				continue
			}

			if err := handler(message); err != nil {
				logger.Error("处理消息失败",
					zap.Error(err),
					zap.String("message_id", message.ID))

				// 重试逻辑
				if message.Retry < 3 {
					message.Retry++
					c.Publish("", queueName, message)
				}

				d.Nack(false, false)
			} else {
				d.Ack(false)
				logger.Info("消息处理成功", zap.String("message_id", message.ID))
			}
		}
	}()

	logger.Info("开始消费消息", zap.String("queue", queueName))
	return nil
}

// DeclareExchange 声明交换器
func (c *Client) DeclareExchange(name, kind string) error {
	return c.channel.ExchangeDeclare(
		name,  // name
		kind,  // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// DeclareQueue 声明队列
func (c *Client) DeclareQueue(name string) error {
	_, err := c.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

// BindQueue 绑定队列到交换器
func (c *Client) BindQueue(queueName, exchangeName, routingKey string) error {
	return c.channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
}

// Close 关闭连接
func (c *Client) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsConnected 检查连接状态
func (c *Client) IsConnected() bool {
	return c.conn != nil && !c.conn.IsClosed()
}

// Reconnect 重新连接
func (c *Client) Reconnect() error {
	c.Close()
	return c.connect()
}

// PublishDelay 发布延迟消息（需要RabbitMQ延迟插件）
func (c *Client) PublishDelay(exchange, routingKey string, message Message, delay time.Duration) error {
	message.Timestamp = time.Now().Unix()

	body, err := json.Marshal(message)
	if err != nil {
		logger.Error("序列化延迟消息失败", zap.Error(err))
		return err
	}

	headers := make(amqp.Table)
	headers["x-delay"] = int32(delay.Milliseconds())

	err = c.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Headers:     headers,
			Body:        body,
		},
	)

	if err != nil {
		logger.Error("发布延迟消息失败", zap.Error(err))
		return err
	}

	logger.Info("延迟消息发布成功",
		zap.String("exchange", exchange),
		zap.String("routing_key", routingKey),
		zap.String("message_id", message.ID),
		zap.Duration("delay", delay))

	return nil
}
