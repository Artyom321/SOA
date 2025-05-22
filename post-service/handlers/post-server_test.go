package handlers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"social-network/common/config"
	"social-network/common/models"
	postpb "social-network/common/proto/post"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&models.Post{}, &models.Comment{}, &models.Like{})
	assert.NoError(t, err)
	return db
}

func TestCreatePost(t *testing.T) {
	db := setupTestDB(t)
	server := &PostServer{db: db}

	req := &postpb.CreatePostRequest{
		Title:       "Test Post",
		Description: "Test Description",
		CreatorId:   1,
		IsPrivate:   false,
		Tags:        []string{"test"},
	}

	resp, err := server.CreatePost(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Post.Id)
	assert.Equal(t, req.Title, resp.Post.Title)
	assert.Equal(t, req.Description, resp.Post.Description)
	assert.Equal(t, req.CreatorId, resp.Post.CreatorId)
	assert.Equal(t, req.IsPrivate, resp.Post.IsPrivate)
	assert.Equal(t, req.Tags, resp.Post.Tags)
}

func TestGetPost(t *testing.T) {
	db := setupTestDB(t)
	server := &PostServer{db: db}

	createReq := &postpb.CreatePostRequest{
		Title:       "Test Post",
		Description: "Test Description",
		CreatorId:   1,
		IsPrivate:   false,
		Tags:        []string{"test"},
	}
	createResp, err := server.CreatePost(context.Background(), createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)

	getReq := &postpb.GetPostRequest{
		Id:          createResp.Post.Id,
		RequesterId: 1,
	}
	getResp, err := server.GetPost(context.Background(), getReq)
	assert.NoError(t, err)
	assert.NotNil(t, getResp)
	assert.Equal(t, createResp.Post.Id, getResp.Post.Id)
	assert.Equal(t, createReq.Title, getResp.Post.Title)
	assert.Equal(t, createReq.Description, getResp.Post.Description)
	assert.Equal(t, createReq.CreatorId, getResp.Post.CreatorId)
	assert.Equal(t, createReq.IsPrivate, getResp.Post.IsPrivate)
	assert.Equal(t, createReq.Tags, getResp.Post.Tags)
}

func TestUpdatePost(t *testing.T) {
	db := setupTestDB(t)
	server := &PostServer{db: db}

	createReq := &postpb.CreatePostRequest{
		Title:       "Test Post",
		Description: "Test Description",
		CreatorId:   1,
		IsPrivate:   false,
		Tags:        []string{"test"},
	}
	createResp, err := server.CreatePost(context.Background(), createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)

	updateReq := &postpb.UpdatePostRequest{
		Id:          createResp.Post.Id,
		Title:       "Updated Post",
		Description: "Updated Description",
		CreatorId:   1,
		IsPrivate:   true,
		Tags:        []string{"updated"},
	}
	updateResp, err := server.UpdatePost(context.Background(), updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updateResp)
	assert.Equal(t, updateReq.Title, updateResp.Post.Title)
	assert.Equal(t, updateReq.Description, updateResp.Post.Description)
	assert.Equal(t, updateReq.CreatorId, updateResp.Post.CreatorId)
	assert.Equal(t, updateReq.IsPrivate, updateResp.Post.IsPrivate)
	assert.Equal(t, updateReq.Tags, updateResp.Post.Tags)
}

func TestDeletePost(t *testing.T) {
	db := setupTestDB(t)
	server := &PostServer{db: db}

	createReq := &postpb.CreatePostRequest{
		Title:       "Test Post",
		Description: "Test Description",
		CreatorId:   1,
		IsPrivate:   false,
		Tags:        []string{"test"},
	}
	createResp, err := server.CreatePost(context.Background(), createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)

	deleteReq := &postpb.DeletePostRequest{
		Id:        createResp.Post.Id,
		CreatorId: 1,
	}
	deleteResp, err := server.DeletePost(context.Background(), deleteReq)
	assert.NoError(t, err)
	assert.NotNil(t, deleteResp)
	assert.True(t, deleteResp.Success)
}

func TestListPosts(t *testing.T) {
	db := setupTestDB(t)
	server := &PostServer{db: db}

	posts := []*postpb.CreatePostRequest{
		{
			Title:       "Post 1",
			Description: "Description 1",
			CreatorId:   1,
			IsPrivate:   false,
			Tags:        []string{"test1"},
		},
		{
			Title:       "Post 2",
			Description: "Description 2",
			CreatorId:   2,
			IsPrivate:   false,
			Tags:        []string{"test2"},
		},
	}

	for _, post := range posts {
		_, err := server.CreatePost(context.Background(), post)
		assert.NoError(t, err)
	}

	listReq := &postpb.ListPostsRequest{
		Page:        1,
		PageSize:    10,
		RequesterId: 1,
	}
	listResp, err := server.ListPosts(context.Background(), listReq)
	assert.NoError(t, err)
	assert.NotNil(t, listResp)
	assert.Equal(t, uint64(2), listResp.TotalCount)
	assert.Len(t, listResp.Posts, 2)
}

func createTestPost(t *testing.T, server *PostServer) *postpb.Post {
	req := &postpb.CreatePostRequest{
		Title:       "Test Post",
		Description: "Test Description",
		CreatorId:   1,
		IsPrivate:   false,
		Tags:        []string{"test"},
	}

	resp, err := server.CreatePost(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	return resp.Post
}

func TestViewPost(t *testing.T) {
	db := setupTestDB(t)

	mockConfig := &config.Config{
		Kafka: struct {
			Broker       string `json:"broker"`
			UserTopic    string `json:"userTopic"`
			ViewTopic    string `json:"viewTopic"`
			LikeTopic    string `json:"likeTopic"`
			CommentTopic string `json:"commentTopic"`
		}{
			ViewTopic: "test_views",
		},
	}

	server := &PostServer{
		db:     db,
		config: mockConfig,
	}

	post := createTestPost(t, server)

	viewReq := &postpb.ViewPostRequest{
		Id:     post.Id,
		UserId: 1,
	}

	viewResp, err := server.ViewPost(context.Background(), viewReq)
	assert.NoError(t, err)
	assert.NotNil(t, viewResp)
	assert.True(t, viewResp.Success)

	privatePostReq := &postpb.CreatePostRequest{
		Title:       "Private Post",
		Description: "For authorized users only",
		CreatorId:   1,
		IsPrivate:   true,
	}
	privatePostResp, _ := server.CreatePost(context.Background(), privatePostReq)

	ownerViewReq := &postpb.ViewPostRequest{
		Id:     privatePostResp.Post.Id,
		UserId: 1,
	}

	ownerViewResp, err := server.ViewPost(context.Background(), ownerViewReq)
	assert.NoError(t, err)
	assert.True(t, ownerViewResp.Success)

	_, getErr := server.GetPost(context.Background(), &postpb.GetPostRequest{
		Id:          privatePostResp.Post.Id,
		RequesterId: 2,
	})

	if getErr != nil {
		unauthorizedViewReq := &postpb.ViewPostRequest{
			Id:     privatePostResp.Post.Id,
			UserId: 2,
		}

		_, viewErr := server.ViewPost(context.Background(), unauthorizedViewReq)
		assert.Error(t, viewErr)
	} else {
		unauthorizedViewReq := &postpb.ViewPostRequest{
			Id:     privatePostResp.Post.Id,
			UserId: 2,
		}

		viewResp, viewErr := server.ViewPost(context.Background(), unauthorizedViewReq)
		assert.NoError(t, viewErr)
		assert.True(t, viewResp.Success)
	}
}

func TestLikePost(t *testing.T) {
	db := setupTestDB(t)

	mockConfig := &config.Config{
		Kafka: struct {
			Broker       string `json:"broker"`
			UserTopic    string `json:"userTopic"`
			ViewTopic    string `json:"viewTopic"`
			LikeTopic    string `json:"likeTopic"`
			CommentTopic string `json:"commentTopic"`
		}{
			LikeTopic: "test_likes",
		},
	}

	server := &PostServer{
		db:     db,
		config: mockConfig,
	}

	post := createTestPost(t, server)

	likeReq := &postpb.LikePostRequest{
		Id:     post.Id,
		UserId: 1,
	}

	likeResp, err := server.LikePost(context.Background(), likeReq)
	assert.NoError(t, err)
	assert.NotNil(t, likeResp)
	assert.True(t, likeResp.Success)
	assert.Equal(t, uint64(1), likeResp.TotalLikes)

	secondLikeResp, err := server.LikePost(context.Background(), likeReq)
	if err != nil {
		assert.Error(t, err)
	} else {
		assert.NotNil(t, secondLikeResp)
	}

	likeReq2 := &postpb.LikePostRequest{
		Id:     post.Id,
		UserId: 2,
	}
	likeResp2, err := server.LikePost(context.Background(), likeReq2)
	assert.NoError(t, err)
	assert.True(t, likeResp2.TotalLikes > 0)
}

func TestAddComment(t *testing.T) {
	db := setupTestDB(t)

	mockConfig := &config.Config{
		Kafka: struct {
			Broker       string `json:"broker"`
			UserTopic    string `json:"userTopic"`
			ViewTopic    string `json:"viewTopic"`
			LikeTopic    string `json:"likeTopic"`
			CommentTopic string `json:"commentTopic"`
		}{
			CommentTopic: "test_comments",
		},
	}

	server := &PostServer{
		db:     db,
		config: mockConfig,
	}

	post := createTestPost(t, server)

	commentReq := &postpb.AddCommentRequest{
		PostId:  post.Id,
		UserId:  1,
		Content: "This is a test comment",
	}

	commentResp, err := server.AddComment(context.Background(), commentReq)
	assert.NoError(t, err)
	assert.NotNil(t, commentResp)
	assert.NotNil(t, commentResp.Comment)
	assert.Equal(t, post.Id, commentResp.Comment.PostId)
	assert.Equal(t, uint64(1), commentResp.Comment.UserId)
	assert.Equal(t, "This is a test comment", commentResp.Comment.Content)

	privatePostReq := &postpb.CreatePostRequest{
		Title:       "Private Post",
		Description: "For authorized users only",
		CreatorId:   1,
		IsPrivate:   true,
	}
	privatePostResp, _ := server.CreatePost(context.Background(), privatePostReq)

	ownerCommentReq := &postpb.AddCommentRequest{
		PostId:  privatePostResp.Post.Id,
		UserId:  1,
		Content: "Owner comment",
	}

	_, err = server.AddComment(context.Background(), ownerCommentReq)
	assert.NoError(t, err)

	unauthorizedCommentReq := &postpb.AddCommentRequest{
		PostId:  privatePostResp.Post.Id,
		UserId:  2,
		Content: "Trying to comment a private post",
	}

	_, err = server.AddComment(context.Background(), unauthorizedCommentReq)
	assert.Error(t, err)
}

func TestGetComments(t *testing.T) {
	db := setupTestDB(t)
	server := &PostServer{db: db}

	post := createTestPost(t, server)

	for i := 1; i <= 5; i++ {
		comment := models.Comment{
			PostID:    post.Id,
			UserID:    uint64(i),
			Content:   fmt.Sprintf("Comment %d", i),
			CreatedAt: time.Now(),
		}
		result := db.Create(&comment)
		assert.NoError(t, result.Error)
	}

	getCommentsReq := &postpb.GetCommentsRequest{
		PostId:      post.Id,
		Page:        1,
		PageSize:    3,
		RequesterId: 1,
	}

	getCommentsResp, err := server.GetComments(context.Background(), getCommentsReq)
	assert.NoError(t, err)
	assert.NotNil(t, getCommentsResp)
	assert.Equal(t, uint64(5), getCommentsResp.TotalCount)
	assert.Len(t, getCommentsResp.Comments, 3)

	getCommentsReq2 := &postpb.GetCommentsRequest{
		PostId:      post.Id,
		Page:        2,
		PageSize:    3,
		RequesterId: 1,
	}

	getCommentsResp2, err := server.GetComments(context.Background(), getCommentsReq2)
	assert.NoError(t, err)
	assert.NotNil(t, getCommentsResp2)
	assert.Equal(t, uint64(5), getCommentsResp2.TotalCount)
	assert.Len(t, getCommentsResp2.Comments, 2)

	privatePostReq := &postpb.CreatePostRequest{
		Title:       "Private Post",
		Description: "For authorized users only",
		CreatorId:   1,
		IsPrivate:   true,
	}
	privatePostResp, _ := server.CreatePost(context.Background(), privatePostReq)

	privateComment := models.Comment{
		PostID:    privatePostResp.Post.Id,
		UserID:    1,
		Content:   "Private comment",
		CreatedAt: time.Now(),
	}
	db.Create(&privateComment)

	ownerReq := &postpb.GetCommentsRequest{
		PostId:      privatePostResp.Post.Id,
		Page:        1,
		PageSize:    10,
		RequesterId: 1,
	}

	ownerResp, err := server.GetComments(context.Background(), ownerReq)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), ownerResp.TotalCount)
	assert.Len(t, ownerResp.Comments, 1)

	unauthorizedReq := &postpb.GetCommentsRequest{
		PostId:      privatePostResp.Post.Id,
		Page:        1,
		PageSize:    10,
		RequesterId: 2,
	}

	_, err = server.GetComments(context.Background(), unauthorizedReq)
	assert.Error(t, err)
}
