package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/IBM/sarama"
	"social-network/common/models"
)

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	db       driver.Conn
	topics   []string
	ready    chan bool
}

func NewKafkaConsumer(brokers []string, groupID string, topics []string, db driver.Conn) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
		db:       db,
		topics:   topics,
		ready:    make(chan bool),
	}, nil
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
	go func() {
		for {
			if err := c.consumer.Consume(ctx, c.topics, c); err != nil {
				log.Printf("Error from consumer: %v", err)
				time.Sleep(time.Second * 5)
			}
			if ctx.Err() != nil {
				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready
	log.Println("Kafka consumer started")
	return nil
}

func (c *KafkaConsumer) Stop() error {
	return c.consumer.Close()
}

func (c *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if err := c.processMessage(message); err != nil {
				log.Printf("Error processing message: %v", err)
			} else {
				session.MarkMessage(message, "")
			}
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *KafkaConsumer) processMessage(message *sarama.ConsumerMessage) error {
	switch message.Topic {
	case "post_views":
		return c.processViewEvent(message.Value)
	case "post_likes":
		return c.processLikeEvent(message.Value)
	case "post_comments":
		return c.processCommentEvent(message.Value)
	default:
		return fmt.Errorf("unknown topic: %s", message.Topic)
	}
}

func (c *KafkaConsumer) processViewEvent(data []byte) error {
	var event models.PostViewEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	query := `
		INSERT INTO post_views (event_time, post_id, user_id)
		VALUES (?, ?, ?)
	`

	return c.db.Exec(context.Background(), query,
		event.ViewedAt,
		event.PostID,
		event.UserID,
	)
}

func (c *KafkaConsumer) processLikeEvent(data []byte) error {
	var event models.PostLikeEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	query := `
		INSERT INTO post_likes (event_time, post_id, user_id)
		VALUES (?, ?, ?)
	`

	return c.db.Exec(context.Background(), query,
		event.LikedAt,
		event.PostID,
		event.UserID,
	)
}

func (c *KafkaConsumer) processCommentEvent(data []byte) error {
	var event models.PostCommentEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	query := `
		INSERT INTO post_comments (event_time, post_id, user_id, comment_id, comment_text)
		VALUES (?, ?, ?, ?, ?)
	`

	return c.db.Exec(context.Background(), query,
		event.CreatedAt,
		event.PostID,
		event.UserID,
		event.CommentID,
		event.CommentText,
	)
}
