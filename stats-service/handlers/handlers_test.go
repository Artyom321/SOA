package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
	statspb "social-network/common/proto/stats"
	"social-network/stats-service/repository"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetPostStats(ctx context.Context, postID uint64) (*repository.PostStats, error) {
	args := m.Called(ctx, postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PostStats), args.Error(1)
}

func (m *MockRepository) GetPostViewsTimeline(ctx context.Context, postID uint64, days uint32) ([]repository.TimelineItem, error) {
	args := m.Called(ctx, postID, days)
	return args.Get(0).([]repository.TimelineItem), args.Error(1)
}

func (m *MockRepository) GetPostLikesTimeline(ctx context.Context, postID uint64, days uint32) ([]repository.TimelineItem, error) {
	args := m.Called(ctx, postID, days)
	return args.Get(0).([]repository.TimelineItem), args.Error(1)
}

func (m *MockRepository) GetPostCommentsTimeline(ctx context.Context, postID uint64, days uint32) ([]repository.TimelineItem, error) {
	args := m.Called(ctx, postID, days)
	return args.Get(0).([]repository.TimelineItem), args.Error(1)
}

func (m *MockRepository) GetTopPosts(ctx context.Context, metricType string, limit int) ([]repository.TopPostItem, error) {
	args := m.Called(ctx, metricType, limit)
	return args.Get(0).([]repository.TopPostItem), args.Error(1)
}

func (m *MockRepository) GetTopUsers(ctx context.Context, metricType string, limit int) ([]repository.TopUserItem, error) {
	args := m.Called(ctx, metricType, limit)
	return args.Get(0).([]repository.TopUserItem), args.Error(1)
}

func setupTestServer() (*StatsServer, *MockRepository) {
	mockRepo := new(MockRepository)
	server := &StatsServer{repo: mockRepo}
	return server, mockRepo
}

func TestGetPostStats(t *testing.T) {
	server, mockRepo := setupTestServer()
	ctx := context.Background()

	postID := uint64(123)

	mockRepo.On("GetPostStats", ctx, postID).Return(&repository.PostStats{
		PostID:        postID,
		ViewsCount:    100,
		LikesCount:    50,
		CommentsCount: 25,
	}, nil)

	req := &statspb.PostStatsRequest{PostId: postID}
	resp, err := server.GetPostStats(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, postID, resp.PostId)
	assert.Equal(t, uint64(100), resp.ViewsCount)
	assert.Equal(t, uint64(50), resp.LikesCount)
	assert.Equal(t, uint64(25), resp.CommentsCount)

	mockRepo.AssertExpectations(t)
}

func TestGetPostViewsTimeline(t *testing.T) {
	server, mockRepo := setupTestServer()
	ctx := context.Background()

	postID := uint64(123)
	days := uint32(7)

	now := time.Now()
	mockTimelineItems := []repository.TimelineItem{
		{Date: now.AddDate(0, 0, -2), Count: 10},
		{Date: now.AddDate(0, 0, -1), Count: 15},
		{Date: now, Count: 20},
	}

	mockRepo.On("GetPostViewsTimeline", ctx, postID, days).Return(mockTimelineItems, nil)

	req := &statspb.PostTimelineRequest{PostId: postID, Days: days}
	resp, err := server.GetPostViewsTimeline(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Items, 3)

	for i, item := range resp.Items {
		assert.Equal(t, mockTimelineItems[i].Count, item.Count)
		assert.Equal(t, timestamppb.New(mockTimelineItems[i].Date).AsTime().Unix(), item.Date.AsTime().Unix())
	}

	mockRepo.AssertExpectations(t)
}

func TestGetPostLikesTimeline(t *testing.T) {
	server, mockRepo := setupTestServer()
	ctx := context.Background()

	postID := uint64(123)
	days := uint32(7)

	now := time.Now()
	mockTimelineItems := []repository.TimelineItem{
		{Date: now.AddDate(0, 0, -2), Count: 5},
		{Date: now.AddDate(0, 0, -1), Count: 10},
		{Date: now, Count: 8},
	}

	mockRepo.On("GetPostLikesTimeline", ctx, postID, days).Return(mockTimelineItems, nil)

	req := &statspb.PostTimelineRequest{PostId: postID, Days: days}
	resp, err := server.GetPostLikesTimeline(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Items, 3)

	for i, item := range resp.Items {
		assert.Equal(t, mockTimelineItems[i].Count, item.Count)
		assert.Equal(t, timestamppb.New(mockTimelineItems[i].Date).AsTime().Unix(), item.Date.AsTime().Unix())
	}

	mockRepo.AssertExpectations(t)
}

func TestGetPostCommentsTimeline(t *testing.T) {
	server, mockRepo := setupTestServer()
	ctx := context.Background()

	postID := uint64(123)
	days := uint32(7)

	now := time.Now()
	mockTimelineItems := []repository.TimelineItem{
		{Date: now.AddDate(0, 0, -2), Count: 3},
		{Date: now.AddDate(0, 0, -1), Count: 7},
		{Date: now, Count: 5},
	}

	mockRepo.On("GetPostCommentsTimeline", ctx, postID, days).Return(mockTimelineItems, nil)

	req := &statspb.PostTimelineRequest{PostId: postID, Days: days}
	resp, err := server.GetPostCommentsTimeline(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Items, 3)

	for i, item := range resp.Items {
		assert.Equal(t, mockTimelineItems[i].Count, item.Count)
		assert.Equal(t, timestamppb.New(mockTimelineItems[i].Date).AsTime().Unix(), item.Date.AsTime().Unix())
	}

	mockRepo.AssertExpectations(t)
}

func TestGetTopPosts(t *testing.T) {
	testCases := []struct {
		name       string
		metricType statspb.MetricType
		metric     string
	}{
		{"Views metric", statspb.MetricType_VIEWS, "views"},
		{"Likes metric", statspb.MetricType_LIKES, "likes"},
		{"Comments metric", statspb.MetricType_COMMENTS, "comments"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server, mockRepo := setupTestServer()
			ctx := context.Background()

			mockTopPosts := []repository.TopPostItem{
				{PostID: 1, Count: 100},
				{PostID: 2, Count: 75},
				{PostID: 3, Count: 50},
			}

			mockRepo.On("GetTopPosts", ctx, tc.metric, 10).Return(mockTopPosts, nil)

			req := &statspb.TopRequest{MetricType: tc.metricType}
			resp, err := server.GetTopPosts(ctx, req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Len(t, resp.Posts, 3)

			for i, post := range resp.Posts {
				assert.Equal(t, mockTopPosts[i].PostID, post.PostId)
				assert.Equal(t, mockTopPosts[i].Count, post.Count)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetTopUsers(t *testing.T) {
	testCases := []struct {
		name       string
		metricType statspb.MetricType
		metric     string
	}{
		{"Views metric", statspb.MetricType_VIEWS, "views"},
		{"Likes metric", statspb.MetricType_LIKES, "likes"},
		{"Comments metric", statspb.MetricType_COMMENTS, "comments"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server, mockRepo := setupTestServer()
			ctx := context.Background()

			mockTopUsers := []repository.TopUserItem{
				{UserID: 1, Count: 100},
				{UserID: 2, Count: 75},
				{UserID: 3, Count: 50},
			}

			mockRepo.On("GetTopUsers", ctx, tc.metric, 10).Return(mockTopUsers, nil)

			req := &statspb.TopRequest{MetricType: tc.metricType}
			resp, err := server.GetTopUsers(ctx, req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Len(t, resp.Users, 3)

			for i, user := range resp.Users {
				assert.Equal(t, mockTopUsers[i].UserID, user.UserId)
				assert.Equal(t, mockTopUsers[i].Count, user.Count)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetPostViewsTimelineEmpty(t *testing.T) {
	server, mockRepo := setupTestServer()
	ctx := context.Background()

	postID := uint64(999)
	days := uint32(7)

	mockRepo.On("GetPostViewsTimeline", ctx, postID, days).Return([]repository.TimelineItem{}, nil)

	req := &statspb.PostTimelineRequest{PostId: postID, Days: days}
	resp, err := server.GetPostViewsTimeline(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.Items)

	mockRepo.AssertExpectations(t)
}
