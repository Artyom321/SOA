package handlers

import (
	"context"
	"errors"
	"strconv"
	"time"

	"social-network/common/models"
	postpb "social-network/common/proto/post"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type PostServer struct {
	postpb.UnimplementedPostServiceServer
	db *gorm.DB
}

func NewPostServer(db *gorm.DB) *postpb.PostServiceServer {
	server := &PostServer{db: db}
	var srv postpb.PostServiceServer = server
	return &srv
}

func (s *PostServer) CreatePost(ctx context.Context, req *postpb.CreatePostRequest) (*postpb.PostResponse, error) {
	creatorID, err := strconv.ParseUint(req.CreatorId, 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid creator ID: %v", err)
	}

	postModel := models.Post{
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   uint(creatorID),
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
	postID, err := strconv.ParseUint(req.Id, 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post ID: %v", err)
	}

	requesterID, err := strconv.ParseUint(req.RequesterId, 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid requester ID: %v", err)
	}

	var postModel models.Post
	if err := s.db.First(&postModel, uint(postID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "post not found")
		}
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if postModel.IsPrivate && postModel.CreatorID != uint(requesterID) {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: this post is private")
	}

	return &postpb.PostResponse{
		Post: convertModelToProto(postModel),
	}, nil
}

func (s *PostServer) UpdatePost(ctx context.Context, req *postpb.UpdatePostRequest) (*postpb.PostResponse, error) {
	postID, err := strconv.ParseUint(req.Id, 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post ID: %v", err)
	}

	creatorID, err := strconv.ParseUint(req.CreatorId, 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid creator ID: %v", err)
	}

	var postModel models.Post
	if err := s.db.First(&postModel, uint(postID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "post not found")
		}
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if postModel.CreatorID != uint(creatorID) {
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
	postID, err := strconv.ParseUint(req.Id, 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post ID: %v", err)
	}

	creatorID, err := strconv.ParseUint(req.CreatorId, 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid creator ID: %v", err)
	}

	var postModel models.Post
	if err := s.db.First(&postModel, uint(postID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "post not found")
		}
		return nil, status.Errorf(codes.Internal, "database error: %v", err)
	}

	if postModel.CreatorID != uint(creatorID) {
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
	requesterID, err := strconv.ParseUint(req.RequesterId, 10, 32)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid requester ID: %v", err)
	}

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
		TotalCount: int32(totalCount),
	}, nil
}

func convertModelToProto(p models.Post) *postpb.Post {
	return &postpb.Post{
		Id:          strconv.FormatUint(uint64(p.ID), 10),
		Title:       p.Title,
		Description: p.Description,
		CreatorId:   strconv.FormatUint(uint64(p.CreatorID), 10),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
		IsPrivate:   p.IsPrivate,
		Tags:        []string(p.Tags),
	}
}
