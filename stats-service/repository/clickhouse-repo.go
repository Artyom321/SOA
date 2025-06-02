package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouseRepo struct {
	conn driver.Conn
}

type PostStats struct {
	PostID        uint64
	ViewsCount    uint64
	LikesCount    uint64
	CommentsCount uint64
}

type TimelineItem struct {
	Date  time.Time
	Count uint64
}

type TopPostItem struct {
	PostID uint64
	Count  uint64
}

type TopUserItem struct {
	UserID uint64
	Count  uint64
}

func NewClickHouseRepo(conn driver.Conn) *ClickHouseRepo {
	return &ClickHouseRepo{conn: conn}
}

func (r *ClickHouseRepo) GetConnection() driver.Conn {
	return r.conn
}

func (r *ClickHouseRepo) GetPostStats(ctx context.Context, postID uint64) (*PostStats, error) {
	stats := &PostStats{
		PostID: postID,
	}

	if err := r.conn.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM post_views WHERE post_id = ?",
		postID,
	).Scan(&stats.ViewsCount); err != nil {
		return nil, fmt.Errorf("failed to get views count: %w", err)
	}

	if err := r.conn.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM post_likes WHERE post_id = ?",
		postID,
	).Scan(&stats.LikesCount); err != nil {
		return nil, fmt.Errorf("failed to get likes count: %w", err)
	}

	if err := r.conn.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM post_comments WHERE post_id = ?",
		postID,
	).Scan(&stats.CommentsCount); err != nil {
		return nil, fmt.Errorf("failed to get comments count: %w", err)
	}

	return stats, nil
}

func (r *ClickHouseRepo) GetPostViewsTimeline(ctx context.Context, postID uint64, days uint32) ([]TimelineItem, error) {
	return r.getPostTimeline(ctx, "post_views", postID, days)
}

func (r *ClickHouseRepo) GetPostLikesTimeline(ctx context.Context, postID uint64, days uint32) ([]TimelineItem, error) {
	return r.getPostTimeline(ctx, "post_likes", postID, days)
}

func (r *ClickHouseRepo) GetPostCommentsTimeline(ctx context.Context, postID uint64, days uint32) ([]TimelineItem, error) {
	return r.getPostTimeline(ctx, "post_comments", postID, days)
}

func (r *ClickHouseRepo) getPostTimeline(ctx context.Context, tableName string, postID uint64, days uint32) ([]TimelineItem, error) {
	var result []TimelineItem

	query := fmt.Sprintf(`
		SELECT 
			date, 
			COUNT(*) as count
		FROM %s
		WHERE 
			post_id = ? AND 
			date >= subtractDays(today(), ?)
		GROUP BY date
		ORDER BY date
	`, tableName)

	rows, err := r.conn.Query(ctx, query, postID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s timeline: %w", tableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var item TimelineItem
		if err := rows.Scan(&item.Date, &item.Count); err != nil {
			return nil, fmt.Errorf("failed to scan timeline row: %w", err)
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating timeline rows: %w", err)
	}

	if len(result) == 0 {
		return []TimelineItem{}, nil
	}

	return result, nil
}

func (r *ClickHouseRepo) GetTopPosts(ctx context.Context, metricType string, limit int) ([]TopPostItem, error) {
	var result []TopPostItem

	var tableName string
	switch metricType {
	case "views":
		tableName = "post_views"
	case "likes":
		tableName = "post_likes"
	case "comments":
		tableName = "post_comments"
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metricType)
	}

	query := fmt.Sprintf(`
		SELECT 
			post_id, 
			COUNT(*) as count
		FROM %s
		GROUP BY post_id
		ORDER BY count DESC
		LIMIT ?
	`, tableName)

	rows, err := r.conn.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item TopPostItem
		if err := rows.Scan(&item.PostID, &item.Count); err != nil {
			return nil, fmt.Errorf("failed to scan top post row: %w", err)
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating top posts rows: %w", err)
	}

	return result, nil
}

func (r *ClickHouseRepo) GetTopUsers(ctx context.Context, metricType string, limit int) ([]TopUserItem, error) {
	var result []TopUserItem

	var tableName string
	switch metricType {
	case "views":
		tableName = "post_views"
	case "likes":
		tableName = "post_likes"
	case "comments":
		tableName = "post_comments"
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metricType)
	}

	query := fmt.Sprintf(`
		SELECT 
			user_id, 
			COUNT(*) as count
		FROM %s
		GROUP BY user_id
		ORDER BY count DESC
		LIMIT ?
	`, tableName)

	rows, err := r.conn.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item TopUserItem
		if err := rows.Scan(&item.UserID, &item.Count); err != nil {
			return nil, fmt.Errorf("failed to scan top user row: %w", err)
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating top users rows: %w", err)
	}

	return result, nil
}

func (r *ClickHouseRepo) InsertViewEvent(ctx context.Context, postID uint64, userID uint64, viewTime time.Time) error {
	query := `
		INSERT INTO post_views (
			event_time, 
			post_id, 
			user_id
		) VALUES (?, ?, ?)
	`

	return r.conn.Exec(ctx, query, viewTime, postID, userID)
}

func (r *ClickHouseRepo) InsertLikeEvent(ctx context.Context, postID uint64, userID uint64, likeTime time.Time) error {
	query := `
		INSERT INTO post_likes (
			event_time, 
			post_id, 
			user_id
		) VALUES (?, ?, ?)
	`

	return r.conn.Exec(ctx, query, likeTime, postID, userID)
}

func (r *ClickHouseRepo) InsertCommentEvent(ctx context.Context, postID uint64, userID uint64, commentID uint64, commentText string, commentTime time.Time) error {
	query := `
		INSERT INTO post_comments (
			event_time, 
			post_id, 
			user_id, 
			comment_id, 
			comment_text
		) VALUES (?, ?, ?, ?, ?)
	`

	return r.conn.Exec(ctx, query, commentTime, postID, userID, commentID, commentText)
}
