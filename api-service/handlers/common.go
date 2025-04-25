package handlers

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"social-network/common/config"
	"social-network/common/kafka"
	postpb "social-network/common/proto/post"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	UserServiceURL    *url.URL
	UserServiceProxy  *httputil.ReverseProxy
	PostServiceClient postpb.PostServiceClient
	Config            *config.Config
	KafkaProducer     *kafka.Producer
}

func NewHandler(userServiceHost string, userServicePort int, postServiceHost string, postServicePort int) *Handler {
	cfg := config.LoadConfig("common/config/config.json")

	userServiceURL, _ := url.Parse(fmt.Sprintf("http://%s:%d", userServiceHost, userServicePort))

	postServiceConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", postServiceHost, postServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to post service: %v", err)
	}

	// Создаем продюсера Kafka, если настроены брокеры
	var kafkaProducer *kafka.Producer
	if len(cfg.Kafka.Brokers) > 0 {
		kafkaProducer = kafka.NewProducer(cfg.Kafka.Brokers)
	}

	h := &Handler{
		UserServiceURL:    userServiceURL,
		PostServiceClient: postpb.NewPostServiceClient(postServiceConn),
		Config:            &cfg,
		KafkaProducer:     kafkaProducer,
	}

	h.UserServiceProxy = h.createUserServiceProxy()

	return h
}
