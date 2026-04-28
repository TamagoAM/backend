package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	PaymentRequestStream = "stream:payment:request"
	PaymentSuccessStream = "stream:payment:success"
	EmailSendStream      = "stream:email:send"
)

// PaymentRequest is the message published to trigger the payment microservice.
type PaymentRequest struct {
	PaymentID int    `json:"payment_id"`
	UserID    int    `json:"user_id"`
	ItemID    int    `json:"item_id"`
	Amount    int    `json:"amount"`
	Currency  string `json:"currency"`
	UserEmail string `json:"user_email"`
	UserName  string `json:"user_name"`
	ItemName  string `json:"item_name"`
}

// PaymentResult is the message received back from the payment microservice.
type PaymentResult struct {
	PaymentID             int    `json:"payment_id"`
	UserID                int    `json:"user_id"`
	ItemID                int    `json:"item_id"`
	Status                string `json:"status"`
	StripePaymentIntentID string `json:"stripe_payment_intent_id,omitempty"`
	ErrorMessage          string `json:"error_message,omitempty"`
}

// EmailRequest is the message published to trigger the email microservice directly.
type EmailRequest struct {
	Template  string            `json:"template"`
	ToAddress string            `json:"to_address"`
	ToName    string            `json:"to_name"`
	Payload   map[string]string `json:"payload"`
}

// RedisStream manages publishing payment requests and consuming results.
type RedisStream struct {
	rdb *redis.Client
}

// NewRedisStream creates a stream client from a Redis URL.
func NewRedisStream(redisURL string) (*RedisStream, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}
	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	// Create consumer group for reading payment results
	_ = rdb.XGroupCreateMkStream(ctx, PaymentSuccessStream, "backend-api", "0").Err()

	return &RedisStream{rdb: rdb}, nil
}

// PublishPaymentRequest publishes a payment request for the payment microservice.
func (rs *RedisStream) PublishPaymentRequest(ctx context.Context, req PaymentRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return rs.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: PaymentRequestStream,
		Values: map[string]interface{}{
			"data": string(data),
		},
	}).Err()
}

// PublishEmailRequest publishes an email request directly to the email microservice stream.
func (rs *RedisStream) PublishEmailRequest(ctx context.Context, req EmailRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return rs.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: EmailSendStream,
		Values: map[string]interface{}{
			"data": string(data),
		},
	}).Err()
}

// ConsumePaymentResults reads payment results from the success stream.
// This should be called in a background goroutine.
func (rs *RedisStream) ConsumePaymentResults(ctx context.Context, handler func(PaymentResult)) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		results, err := rs.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "backend-api",
			Consumer: "backend-api-1",
			Streams:  []string{PaymentSuccessStream, ">"},
			Count:    10,
			Block:    5 * time.Second,
		}).Result()
		if err != nil {
			if err == redis.Nil || ctx.Err() != nil {
				continue
			}
			log.Printf("[store-stream] read error: %v", err)
			continue
		}

		for _, stream := range results {
			for _, msg := range stream.Messages {
				dataStr, ok := msg.Values["data"].(string)
				if !ok {
					continue
				}
				var result PaymentResult
				if err := json.Unmarshal([]byte(dataStr), &result); err != nil {
					log.Printf("[store-stream] unmarshal error: %v", err)
					_ = rs.rdb.XAck(ctx, PaymentSuccessStream, "backend-api", msg.ID).Err()
					continue
				}

				handler(result)

				_ = rs.rdb.XAck(ctx, PaymentSuccessStream, "backend-api", msg.ID).Err()
			}
		}
	}
}

// Close closes the Redis connection.
func (rs *RedisStream) Close() error {
	return rs.rdb.Close()
}
