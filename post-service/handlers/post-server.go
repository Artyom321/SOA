package handlers

import (
	"context"
	"errors"
	"log"
	"time"

	"social-network/common/config"
	"social-network/common/kafka"
	"social-network/common/models"
	postpb "social-network/common/proto/post"

	"github.com/IBM/sarama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type PostServer struct {
	postpb.UnimplementedPostServiceServer
	db            *gorm.DB
	kafkaProducer sarama.SyncProducer
	config        *config.Config
}

func NewPostServer(db *gorm.DB, config *config.Config) *postpb.PostServiceServer {
	producer := kafka.NewProducer(config.Kafka.Broker)

	server := &PostServer{
		db:            db,
		kafkaProducer: producer,
		config:        config,
	}
	var srv postpb.PostServiceServer = server
	return &srv
}

func (s *PostServer) CreatePost(ctx context.Context, req *postpb.CreatePostRequest) (*postpb.PostResponse, error) {
	postModel := models.Post{
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   req.CreatorId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsPrivate:   req.IsPrivate,
		Tags:        models.StringArray(req.Tags),
	}

	if err := s.db.Create(&postModel).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create post: %v", err)
	}

	return &postpb.PostResponse{
		Post: convertModelToProto(postModel),
	}, nil
}

func (s *PostServer) GetPost(ctx context.Context, req *postpb.GetPostRequest) (*postpb.PostResponse, error) {
	postID := req.Id

	requesterID := req.RequesterId

	var postModel models.Post
	if err := s.db.First(&postModel, uint(postID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "post not found")
		}
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if postModel.IsPrivate && postModel.CreatorID != requesterID {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: this post is private")
	}

	return &postpb.PostResponse{
		Post: convertModelToProto(postModel),
	}, nil
}

func (s *PostServer) UpdatePost(ctx context.Context, req *postpb.UpdatePostRequest) (*postpb.PostResponse, error) {
	postID := req.Id

	creatorID := req.CreatorId

	var postModel models.Post
	if err := s.db.First(&postModel, uint(postID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "post not found")
		}
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if postModel.CreatorID != creatorID {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: you can only update your own posts")
	}

	if req.Title != "" {
		postModel.Title = req.Title
	}
	if req.Description != "" {
		postModel.Description = req.Description
	}
	postModel.IsPrivate = req.IsPrivate
	if len(req.Tags) > 0 {
		postModel.Tags = models.StringArray(req.Tags)
	}
	postModel.UpdatedAt = time.Now()

	if err := s.db.Save(&postModel).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update post: %v", err)
	}

	return &postpb.PostResponse{
		Post: convertModelToProto(postModel),
	}, nil
}

func (s *PostServer) DeletePost(ctx context.Context, req *postpb.DeletePostRequest) (*postpb.DeletePostResponse, error) {
	postID := req.Id

	creatorID := req.CreatorId

	var postModel models.Post
	if err := s.db.First(&postModel, uint(postID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "post not found")
		}
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if postModel.CreatorID != creatorID {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: you can only delete your own posts")
	}

	if err := s.db.Delete(&postModel).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete post: %v", err)
	}

	return &postpb.DeletePostResponse{
		Success: true,
	}, nil
}

func (s *PostServer) ListPosts(ctx context.Context, req *postpb.ListPostsRequest) (*postpb.ListPostsResponse, error) {
	requesterID := req.RequesterId

	offset := (req.Page - 1) * req.PageSize
	limit := req.PageSize

	var posts []models.Post
	var totalCount int64

	query := s.db.Model(&models.Post{}).Where("is_private = false OR creator_id = ?", uint(requesterID))
	query.Count(&totalCount)

	if err := query.Offset(int(offset)).Limit(int(limit)).Find(&posts).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list posts: %v", err)
	}

	postProtos := make([]*postpb.Post, len(posts))
	for i, post := range posts {
		postProtos[i] = convertModelToProto(post)
	}

	return &postpb.ListPostsResponse{
		Posts:      postProtos,
		TotalCount: uint64(totalCount),
	}, nil
}

func convertModelToProto(p models.Post) *postpb.Post {
	return &postpb.Post{
		Id:          uint64(p.ID),
		Title:       p.Title,
		Description: p.Description,
		CreatorId:   uint64(p.CreatorID),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
		IsPrivate:   p.IsPrivate,
		Tags:        []string(p.Tags),
	}
}

func (s *PostServer) ViewPost(ctx context.Context, req *postpb.ViewPostRequest) (*postpb.ViewPostResponse, error) {
	var post models.Post
	if err := s.db.First(&post, req.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "post not found")
		}
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if post.IsPrivate && post.CreatorID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: this post is private")
	}

	if s.kafkaProducer != nil {
		event := models.PostViewEvent{
			UserID:   req.UserId,
			PostID:   req.Id,
			ViewedAt: time.Now(),
		}

		err := kafka.SendMessage(s.kafkaProducer, s.config.Kafka.ViewTopic, event)
		if err != nil {
			log.Printf("Error sending view event to Kafka: %v", err)
		}
	}

	return &postpb.ViewPostResponse{
		Success: true,
	}, nil
}

func (s *PostServer) LikePost(ctx context.Context, req *postpb.LikePostRequest) (*postpb.LikePostResponse, error) {
	var post models.Post
	if err := s.db.First(&post, req.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "post not found")
		}
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if post.IsPrivate && post.CreatorID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: this post is private")
	}

	var existingLike models.Like
	result := s.db.Where("post_id = ? AND user_id = ?", req.Id, req.UserId).First(&existingLike)

	if result.Error == nil {
		return nil, status.Errorf(codes.AlreadyExists, "you have already liked this post")
	}

	like := models.Like{
		PostID:    req.Id,
		UserID:    req.UserId,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(&like).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save like: %v", err)
	}

	if s.kafkaProducer != nil {
		event := models.PostLikeEvent{
			UserID:  req.UserId,
			PostID:  req.Id,
			LikedAt: time.Now(),
		}

		err := kafka.SendMessage(s.kafkaProducer, s.config.Kafka.LikeTopic, event)
		if err != nil {
			log.Printf("Error sending like event to Kafka: %v", err)
		}
	}

	var totalLikes int64
	s.db.Model(&models.Like{}).Where("post_id = ?", req.Id).Count(&totalLikes)

	return &postpb.LikePostResponse{
		Success:    true,
		TotalLikes: uint64(totalLikes),
	}, nil
}

func (s *PostServer) AddComment(ctx context.Context, req *postpb.AddCommentRequest) (*postpb.CommentResponse, error) {
	var post models.Post
	if err := s.db.First(&post, req.PostId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "post not found")
		}
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if post.IsPrivate && post.CreatorID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: this post is private")
	}

	comment := models.Comment{
		PostID:    req.PostId,
		UserID:    req.UserId,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(&comment).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save comment: %v", err)
	}

	if s.kafkaProducer != nil {
		event := models.PostCommentEvent{
			UserID:      req.UserId,
			PostID:      req.PostId,
			CommentID:   comment.ID,
			CommentText: comment.Content,
			CreatedAt:   comment.CreatedAt,
		}

		err := kafka.SendMessage(s.kafkaProducer, s.config.Kafka.CommentTopic, event)
		if err != nil {
			log.Printf("Error sending comment event to Kafka: %v", err)
		}
	}

	return &postpb.CommentResponse{
		Comment: &postpb.Comment{
			Id:        comment.ID,
			PostId:    comment.PostID,
			UserId:    comment.UserID,
			Content:   comment.Content,
			CreatedAt: timestamppb.New(comment.CreatedAt),
		},
	}, nil
}

func (s *PostServer) GetComments(ctx context.Context, req *postpb.GetCommentsRequest) (*postpb.GetCommentsResponse, error) {
	var post models.Post
	if err := s.db.First(&post, req.PostId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "post not found")
		}
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if post.IsPrivate && post.CreatorID != req.RequesterId {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: this post is private")
	}

	page := int(req.Page)
	if page < 1 {
		page = 1
	}

	pageSize := int(req.PageSize)
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var comments []models.Comment
	var totalCount int64

	if err := s.db.Model(&models.Comment{}).Where("post_id = ?", req.PostId).Count(&totalCount).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count comments: %v", err)
	}

	if err := s.db.Where("post_id = ?", req.PostId).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve comments: %v", err)
	}

	commentProtos := make([]*postpb.Comment, len(comments))
	for i, comment := range comments {
		commentProtos[i] = &postpb.Comment{
			Id:        comment.ID,
			PostId:    comment.PostID,
			UserId:    comment.UserID,
			Content:   comment.Content,
			CreatedAt: timestamppb.New(comment.CreatedAt),
		}
	}

	return &postpb.GetCommentsResponse{
		Comments:   commentProtos,
		TotalCount: uint64(totalCount),
	}, nil
}
