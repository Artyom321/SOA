package handlers

import (
	"context"
	"testing"

	"social-network/common/models"
	postpb "social-network/common/proto/post"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&models.Post{})
	assert.NoError(t, err)
	return db
}

func TestCreatePost(t *testing.T) {
	db := setupTestDB(t)
	server := &PostServer{db: db}

	req := &postpb.CreatePostRequest{
		Title:       "Test Post",
		Description: "Test Description",
		CreatorId:   "1",
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

	// Create a post first
	createReq := &postpb.CreatePostRequest{
		Title:       "Test Post",
		Description: "Test Description",
		CreatorId:   "1",
		IsPrivate:   false,
		Tags:        []string{"test"},
	}
	createResp, err := server.CreatePost(context.Background(), createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)

	// Get the created post
	getReq := &postpb.GetPostRequest{
		Id:          createResp.Post.Id,
		RequesterId: "1",
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

	// Create a post first
	createReq := &postpb.CreatePostRequest{
		Title:       "Test Post",
		Description: "Test Description",
		CreatorId:   "1",
		IsPrivate:   false,
		Tags:        []string{"test"},
	}
	createResp, err := server.CreatePost(context.Background(), createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)

	// Update the post
	updateReq := &postpb.UpdatePostRequest{
		Id:          createResp.Post.Id,
		Title:       "Updated Post",
		Description: "Updated Description",
		CreatorId:   "1",
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

	// Create a post first
	createReq := &postpb.CreatePostRequest{
		Title:       "Test Post",
		Description: "Test Description",
		CreatorId:   "1",
		IsPrivate:   false,
		Tags:        []string{"test"},
	}
	createResp, err := server.CreatePost(context.Background(), createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)

	// Delete the post
	deleteReq := &postpb.DeletePostRequest{
		Id:        createResp.Post.Id,
		CreatorId: "1",
	}
	deleteResp, err := server.DeletePost(context.Background(), deleteReq)
	assert.NoError(t, err)
	assert.NotNil(t, deleteResp)
	assert.True(t, deleteResp.Success)
}

func TestListPosts(t *testing.T) {
	db := setupTestDB(t)
	server := &PostServer{db: db}

	// Create multiple posts
	posts := []*postpb.CreatePostRequest{
		{
			Title:       "Post 1",
			Description: "Description 1",
			CreatorId:   "1",
			IsPrivate:   false,
			Tags:        []string{"test1"},
		},
		{
			Title:       "Post 2",
			Description: "Description 2",
			CreatorId:   "2",
			IsPrivate:   false,
			Tags:        []string{"test2"},
		},
	}

	for _, post := range posts {
		_, err := server.CreatePost(context.Background(), post)
		assert.NoError(t, err)
	}

	// List posts
	listReq := &postpb.ListPostsRequest{
		Page:        1,
		PageSize:    10,
		RequesterId: "1",
	}
	listResp, err := server.ListPosts(context.Background(), listReq)
	assert.NoError(t, err)
	assert.NotNil(t, listResp)
	assert.Equal(t, int32(2), listResp.TotalCount)
	assert.Len(t, listResp.Posts, 2)
}
