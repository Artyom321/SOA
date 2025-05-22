package handlers

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"social-network/common/config"
	"social-network/common/kafka"
	postpb "social-network/common/proto/post"

	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	UserServiceURL    *url.URL
	UserServiceProxy  *httputil.ReverseProxy
	PostServiceClient postpb.PostServiceClient
	Config            *config.Config
	KafkaProducer     sarama.SyncProducer
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

	var kafkaProducer sarama.SyncProducer
	if len(cfg.Kafka.Broker) > 0 {
		kafkaProducer = kafka.NewProducer(cfg.Kafka.Broker)
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
