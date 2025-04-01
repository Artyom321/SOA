package handlers

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	postpb "social-network/common/proto/post"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	UserServiceURL    *url.URL
	UserServiceProxy  *httputil.ReverseProxy
	PostServiceClient postpb.PostServiceClient
}

func NewHandler(userServiceHost string, userServicePort int, postServiceHost string, postServicePort int) *Handler {
	userServiceURL, _ := url.Parse(fmt.Sprintf("http://%s:%d", userServiceHost, userServicePort))

	postServiceConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", postServiceHost, postServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to post service: %v", err)
	}

	h := &Handler{
		UserServiceURL:    userServiceURL,
		PostServiceClient: postpb.NewPostServiceClient(postServiceConn),
	}

	h.UserServiceProxy = h.createUserServiceProxy()

	return h
}
