package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	statspb "social-network/common/proto/stats"
	"social-network/stats-service/repository"
)

type StatsServer struct {
	statspb.UnimplementedStatsServiceServer
	repo repository.StatsRepository
}

func NewStatsServer(repo repository.StatsRepository) *statspb.StatsServiceServer {
	server := &StatsServer{repo: repo}
	var srv statspb.StatsServiceServer = server
	return &srv
}

func (s *StatsServer) GetPostStats(ctx context.Context, req *statspb.PostStatsRequest) (*statspb.PostStatsResponse, error) {
	stats, err := s.repo.GetPostStats(ctx, req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post stats: %v", err)
	}

	return &statspb.PostStatsResponse{
		PostId:        stats.PostID,
		ViewsCount:    stats.ViewsCount,
		LikesCount:    stats.LikesCount,
		CommentsCount: stats.CommentsCount,
	}, nil
}

func (s *StatsServer) GetPostViewsTimeline(ctx context.Context, req *statspb.PostTimelineRequest) (*statspb.TimelineResponse, error) {
	timelineItems, err := s.repo.GetPostViewsTimeline(ctx, req.PostId, req.Days)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post views timeline: %v", err)
	}

	return convertTimelineResponse(timelineItems), nil
}

func (s *StatsServer) GetPostLikesTimeline(ctx context.Context, req *statspb.PostTimelineRequest) (*statspb.TimelineResponse, error) {
	timelineItems, err := s.repo.GetPostLikesTimeline(ctx, req.PostId, req.Days)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post likes timeline: %v", err)
	}

	return convertTimelineResponse(timelineItems), nil
}

func (s *StatsServer) GetPostCommentsTimeline(ctx context.Context, req *statspb.PostTimelineRequest) (*statspb.TimelineResponse, error) {
	timelineItems, err := s.repo.GetPostCommentsTimeline(ctx, req.PostId, req.Days)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post comments timeline: %v", err)
	}

	return convertTimelineResponse(timelineItems), nil
}

func (s *StatsServer) GetTopPosts(ctx context.Context, req *statspb.TopRequest) (*statspb.TopPostsResponse, error) {
	metricType := getMetricTypeName(req.MetricType)

	topPosts, err := s.repo.GetTopPosts(ctx, metricType, 10)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get top posts: %v", err)
	}

	response := &statspb.TopPostsResponse{
		Posts: make([]*statspb.TopPostItem, len(topPosts)),
	}

	for i, post := range topPosts {
		response.Posts[i] = &statspb.TopPostItem{
			PostId: post.PostID,
			Count:  post.Count,
		}
	}

	return response, nil
}

func (s *StatsServer) GetTopUsers(ctx context.Context, req *statspb.TopRequest) (*statspb.TopUsersResponse, error) {
	metricType := getMetricTypeName(req.MetricType)

	topUsers, err := s.repo.GetTopUsers(ctx, metricType, 10)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get top users: %v", err)
	}

	response := &statspb.TopUsersResponse{
		Users: make([]*statspb.TopUserItem, len(topUsers)),
	}

	for i, user := range topUsers {
		response.Users[i] = &statspb.TopUserItem{
			UserId: user.UserID,
			Count:  user.Count,
		}
	}

	return response, nil
}

func getMetricTypeName(metricType statspb.MetricType) string {
	switch metricType {
	case statspb.MetricType_VIEWS:
		return "views"
	case statspb.MetricType_LIKES:
		return "likes"
	case statspb.MetricType_COMMENTS:
		return "comments"
	default:
		return "views"
	}
}

func convertTimelineResponse(items []repository.TimelineItem) *statspb.TimelineResponse {
	response := &statspb.TimelineResponse{
		Items: make([]*statspb.TimelineItem, len(items)),
	}

	for i, item := range items {
		response.Items[i] = &statspb.TimelineItem{
			Date:  timestamppb.New(item.Date),
			Count: item.Count,
		}
	}

	return response
}
