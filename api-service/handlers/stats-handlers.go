package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"social-network/common/models"
	statspb "social-network/common/proto/stats"
)

// GetPostStatsHandler обрабатывает запросы на получение статистики по посту
// @Summary Получение статистики по посту
// @Description Возвращает количество просмотров, лайков и комментариев по посту
// @Tags Statistics
// @Produce json
// @Param id path string true "ID поста"
// @Success 200 {object} models.PostStatsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /stats/posts/{id} [get]
func (h *Handler) GetPostStatsHandler(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid post ID"})
		return
	}

	req := &statspb.PostStatsRequest{
		PostId: postID,
	}

	resp, err := h.StatsServiceClient.GetPostStats(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get post stats"})
		return
	}

	c.JSON(http.StatusOK, models.PostStatsResponse{
		PostID:        resp.PostId,
		ViewsCount:    resp.ViewsCount,
		LikesCount:    resp.LikesCount,
		CommentsCount: resp.CommentsCount,
	})
}

// GetPostViewsTimelineHandler обрабатывает запросы на получение динамики просмотров
// @Summary Получение динамики просмотров по посту
// @Description Возвращает массив с датой (день) и количеством просмотров за день
// @Tags Statistics
// @Produce json
// @Param id path string true "ID поста"
// @Param days query int false "Количество дней (по умолчанию 30)"
// @Success 200 {object} models.TimelineResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /stats/posts/{id}/views/timeline [get]
func (h *Handler) GetPostViewsTimelineHandler(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid post ID"})
		return
	}

	days, err := strconv.ParseUint(c.DefaultQuery("days", "30"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid days parameter"})
		return
	}

	req := &statspb.PostTimelineRequest{
		PostId: postID,
		Days:   uint32(days),
	}

	resp, err := h.StatsServiceClient.GetPostViewsTimeline(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get views timeline"})
		return
	}

	c.JSON(http.StatusOK, convertTimelineResponse(resp))
}

// GetPostLikesTimelineHandler обрабатывает запросы на получение динамики лайков
// @Summary Получение динамики лайков по посту
// @Description Возвращает массив с датой (день) и количеством лайков за день
// @Tags Statistics
// @Produce json
// @Param id path string true "ID поста"
// @Param days query int false "Количество дней (по умолчанию 30)"
// @Success 200 {object} models.TimelineResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /stats/posts/{id}/likes/timeline [get]
func (h *Handler) GetPostLikesTimelineHandler(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid post ID"})
		return
	}

	days, err := strconv.ParseUint(c.DefaultQuery("days", "30"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid days parameter"})
		return
	}

	req := &statspb.PostTimelineRequest{
		PostId: postID,
		Days:   uint32(days),
	}

	resp, err := h.StatsServiceClient.GetPostLikesTimeline(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get likes timeline"})
		return
	}

	c.JSON(http.StatusOK, convertTimelineResponse(resp))
}

// GetPostCommentsTimelineHandler обрабатывает запросы на получение динамики комментариев
// @Summary Получение динамики комментариев по посту
// @Description Возвращает массив с датой (день) и количеством комментариев за день
// @Tags Statistics
// @Produce json
// @Param id path string true "ID поста"
// @Param days query int false "Количество дней (по умолчанию 30)"
// @Success 200 {object} models.TimelineResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /stats/posts/{id}/comments/timeline [get]
func (h *Handler) GetPostCommentsTimelineHandler(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid post ID"})
		return
	}

	days, err := strconv.ParseUint(c.DefaultQuery("days", "30"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid days parameter"})
		return
	}

	req := &statspb.PostTimelineRequest{
		PostId: postID,
		Days:   uint32(days),
	}

	resp, err := h.StatsServiceClient.GetPostCommentsTimeline(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get comments timeline"})
		return
	}

	c.JSON(http.StatusOK, convertTimelineResponse(resp))
}

// GetTopPostsHandler обрабатывает запросы на получение топ-10 постов
// @Summary Получение топ-10 постов
// @Description Возвращает топ-10 постов по количеству просмотров, лайков или комментариев
// @Tags Statistics
// @Produce json
// @Param metric_type query string true "Тип метрики (views, likes, comments)"
// @Success 200 {object} models.TopPostsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /stats/posts/top [get]
func (h *Handler) GetTopPostsHandler(c *gin.Context) {
	metricType := c.Query("metric_type")
	var pbMetricType statspb.MetricType

	switch metricType {
	case "views":
		pbMetricType = statspb.MetricType_VIEWS
	case "likes":
		pbMetricType = statspb.MetricType_LIKES
	case "comments":
		pbMetricType = statspb.MetricType_COMMENTS
	default:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid metric_type. Must be one of: views, likes, comments"})
		return
	}

	req := &statspb.TopRequest{
		MetricType: pbMetricType,
	}

	resp, err := h.StatsServiceClient.GetTopPosts(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get top posts"})
		return
	}

	result := models.TopPostsResponse{
		Posts: make([]models.TopPostItem, len(resp.Posts)),
	}

	for i, post := range resp.Posts {
		result.Posts[i] = models.TopPostItem{
			PostID: post.PostId,
			Count:  post.Count,
		}
	}

	c.JSON(http.StatusOK, result)
}

// GetTopUsersHandler обрабатывает запросы на получение топ-10 пользователей
// @Summary Получение топ-10 пользователей
// @Description Возвращает топ-10 пользователей по количеству просмотров, лайков или комментариев
// @Tags Statistics
// @Produce json
// @Param metric_type query string true "Тип метрики (views, likes, comments)"
// @Success 200 {object} models.TopUsersResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /stats/users/top [get]
func (h *Handler) GetTopUsersHandler(c *gin.Context) {
	metricType := c.Query("metric_type")
	var pbMetricType statspb.MetricType

	switch metricType {
	case "views":
		pbMetricType = statspb.MetricType_VIEWS
	case "likes":
		pbMetricType = statspb.MetricType_LIKES
	case "comments":
		pbMetricType = statspb.MetricType_COMMENTS
	default:
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid metric_type. Must be one of: views, likes, comments"})
		return
	}

	req := &statspb.TopRequest{
		MetricType: pbMetricType,
	}

	resp, err := h.StatsServiceClient.GetTopUsers(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get top users"})
		return
	}

	result := models.TopUsersResponse{
		Users: make([]models.TopUserItem, len(resp.Users)),
	}

	for i, user := range resp.Users {
		result.Users[i] = models.TopUserItem{
			UserID: user.UserId,
			Count:  user.Count,
		}
	}

	c.JSON(http.StatusOK, result)
}

// Вспомогательная функция для конвертации ответа временной шкалы
func convertTimelineResponse(resp *statspb.TimelineResponse) models.TimelineResponse {
	result := models.TimelineResponse{
		Items: make([]models.TimelineItem, len(resp.Items)),
	}

	for i, item := range resp.Items {
		result.Items[i] = models.TimelineItem{
			Date:  item.Date.AsTime(),
			Count: item.Count,
		}
	}

	return result
}
