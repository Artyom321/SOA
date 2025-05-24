package repository

import (
	"context"
)

type StatsRepository interface {
	GetPostStats(ctx context.Context, postID uint64) (*PostStats, error)
	GetPostViewsTimeline(ctx context.Context, postID uint64, days uint32) ([]TimelineItem, error)
	GetPostLikesTimeline(ctx context.Context, postID uint64, days uint32) ([]TimelineItem, error)
	GetPostCommentsTimeline(ctx context.Context, postID uint64, days uint32) ([]TimelineItem, error)
	GetTopPosts(ctx context.Context, metricType string, limit int) ([]TopPostItem, error)
	GetTopUsers(ctx context.Context, metricType string, limit int) ([]TopUserItem, error)
}
