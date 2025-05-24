package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"social-network/common/config"
	statspb "social-network/common/proto/stats"
	"social-network/stats-service/consumer"
	"social-network/stats-service/handlers"
	"social-network/stats-service/repository"
)

func main() {
	serviceConfig := config.LoadConfig("common/config/config.json")

	clickhouseHost := getEnv("CLICKHOUSE_HOST", "clickhouse")
	clickhousePortStr := getEnv("CLICKHOUSE_PORT", "9000")
	clickhousePort, err := strconv.Atoi(clickhousePortStr)
	if err != nil {
		clickhousePort = 9000
		log.Printf("Invalid CLICKHOUSE_PORT, using default: %d", clickhousePort)
	}

	clickhouseUser := getEnv("CLICKHOUSE_USER", "default")
	clickhousePassword := getEnv("CLICKHOUSE_PASSWORD", "password")
	clickhouseDB := getEnv("CLICKHOUSE_DB", "social_stats")

	log.Printf("Connecting to ClickHouse at %s:%d...", clickhouseHost, clickhousePort)

	var conn driver.Conn
	maxRetries := 30
	retryDelay := time.Second * 5

	for retry := 0; retry < maxRetries; retry++ {
		log.Printf("Attempt %d/%d to connect to ClickHouse...", retry+1, maxRetries)

		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%d", clickhouseHost, clickhousePort)},
			Auth: clickhouse.Auth{
				Database: clickhouseDB,
				Username: clickhouseUser,
				Password: clickhousePassword,
			},
			Settings: clickhouse.Settings{
				"max_execution_time": 60,
			},
			DialTimeout:     time.Second * 10,
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour,
		})

		if err == nil {
			if err = conn.Ping(context.Background()); err == nil {
				log.Println("Successfully connected to ClickHouse!")
				break
			}
			log.Printf("Failed to ping ClickHouse: %v", err)
		} else {
			log.Printf("Failed to open connection to ClickHouse: %v", err)
		}

		if retry < maxRetries-1 {
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		} else {
			log.Fatalf("Failed to connect to ClickHouse after %d attempts", maxRetries)
		}
	}

	if err := initClickHouseTables(conn); err != nil {
		log.Fatalf("Failed to initialize ClickHouse tables: %v", err)
	}

	repo := repository.NewClickHouseRepo(conn)

	statsServer := handlers.NewStatsServer(repo)
	grpcServer := grpc.NewServer()
	statspb.RegisterStatsServiceServer(grpcServer, *statsServer)
	reflection.Register(grpcServer)

	port := serviceConfig.StatsService.Port
	if port == 0 {
		port = 8083
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	kafkaBrokers := getEnv("KAFKA_BROKERS", "kafka:29092")
	kafkaGroupID := "stats-consumer"
	kafkaTopics := []string{
		getEnv("KAFKA_VIEW_TOPIC", "post_views"),
		getEnv("KAFKA_LIKE_TOPIC", "post_likes"),
		getEnv("KAFKA_COMMENT_TOPIC", "post_comments"),
	}

	consumer, err := consumer.NewKafkaConsumer(
		strings.Split(kafkaBrokers, ","),
		kafkaGroupID,
		kafkaTopics,
		conn,
	)
	if err != nil {
		log.Printf("Warning: Failed to create Kafka consumer: %v", err)
		log.Println("Stats service will run without Kafka consumer.")
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := consumer.Start(ctx); err != nil {
			log.Printf("Warning: Failed to start Kafka consumer: %v", err)
			log.Println("Stats service will run without Kafka consumer.")
		} else {
			defer consumer.Stop()
			log.Println("Kafka consumer started successfully")
		}
	}

	log.Printf("Stats-service gRPC server running on port %d\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func initClickHouseTables(conn driver.Conn) error {
	log.Println("Initializing ClickHouse tables...")

	queries := []string{
		`CREATE TABLE IF NOT EXISTS post_views (
			event_time DateTime,
			date Date DEFAULT toDate(event_time),
			post_id UInt64,
			user_id UInt64
		) ENGINE = MergeTree()
		ORDER BY (date, post_id, user_id)`,

		`CREATE TABLE IF NOT EXISTS post_likes (
			event_time DateTime,
			date Date DEFAULT toDate(event_time),
			post_id UInt64,
			user_id UInt64
		) ENGINE = MergeTree()
		ORDER BY (date, post_id, user_id)`,

		`CREATE TABLE IF NOT EXISTS post_comments (
			event_time DateTime,
			date Date DEFAULT toDate(event_time),
			post_id UInt64,
			user_id UInt64,
			comment_id UInt64,
			comment_text String
		) ENGINE = MergeTree()
		ORDER BY (date, post_id, user_id)`,
	}

	for _, query := range queries {
		err := conn.Exec(context.Background(), query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %v", query, err)
		}
	}

	log.Println("ClickHouse tables initialized successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
