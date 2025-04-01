package handlers

import (
	"log"
	"net/http"
	"social-network/common/models"
	"strconv"

	postpb "social-network/common/proto/post"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler обрабатывает запросы на создание поста
// @Summary Создание нового поста
// @Description Создает новый пост от имени текущего пользователя
// @Tags Posts
// @Accept json
// @Produce json
// @Security sessionAuth
// @Param body body models.CreatePostRequest true "Данные для создания поста"
// @Success 201 {object} models.PostResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts [post]
func (h *Handler) CreatePostHandler(c *gin.Context) {
	userID, exists := c.Get(UserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User ID not found"})
		return
	}

	var reqBody models.CreatePostRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	req := &postpb.CreatePostRequest{
		Title:       reqBody.Title,
		Description: reqBody.Description,
		CreatorId:   userID.(string),
		IsPrivate:   reqBody.IsPrivate,
		Tags:        reqBody.Tags,
	}

	resp, err := h.PostServiceClient.CreatePost(c, req)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create post"})
		return
	}

	post := resp.GetPost()
	if post == nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Invalid response from post service"})
		return
	}

	c.JSON(http.StatusCreated, convertProtoToResponse(post))
}

// GetPostHandler обрабатывает запросы на получение поста по ID
// @Summary Получение поста по ID
// @Description Возвращает пост по его ID
// @Tags Posts
// @Produce json
// @Security sessionAuth
// @Param id path string true "ID поста"
// @Success 200 {object} models.PostResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /posts/{id} [get]
func (h *Handler) GetPostHandler(c *gin.Context) {
	userID, exists := c.Get(UserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User ID not found"})
		return
	}

	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Post ID is required"})
		return
	}

	req := &postpb.GetPostRequest{
		Id:          postID,
		RequesterId: userID.(string),
	}

	resp, err := h.PostServiceClient.GetPost(c, req)
	if err != nil {
		log.Printf("Error getting post: %v", err)
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Post not found or access denied"})
		return
	}

	post := resp.GetPost()
	if post == nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Invalid response from post service"})
		return
	}

	c.JSON(http.StatusOK, convertProtoToResponse(post))
}

// UpdatePostHandler обрабатывает запросы на обновление поста
// @Summary Обновление поста
// @Description Обновляет данные поста
// @Tags Posts
// @Accept json
// @Produce json
// @Security sessionAuth
// @Param id path string true "ID поста"
// @Param body body models.UpdatePostRequest true "Новые данные для поста"
// @Success 200 {object} models.PostResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [put]
func (h *Handler) UpdatePostHandler(c *gin.Context) {
	userID, exists := c.Get(UserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User ID not found"})
		return
	}

	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Post ID is required"})
		return
	}

	var reqBody models.UpdatePostRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	req := &postpb.UpdatePostRequest{
		Id:          postID,
		Title:       reqBody.Title,
		Description: reqBody.Description,
		CreatorId:   userID.(string),
		IsPrivate:   reqBody.IsPrivate,
		Tags:        reqBody.Tags,
	}

	resp, err := h.PostServiceClient.UpdatePost(c, req)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update post"})
		return
	}

	post := resp.GetPost()
	if post == nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Invalid response from post service"})
		return
	}

	c.JSON(http.StatusOK, convertProtoToResponse(post))
}

// DeletePostHandler обрабатывает запросы на удаление поста
// @Summary Удаление поста
// @Description Удаляет пост по его ID
// @Tags Posts
// @Security sessionAuth
// @Param id path string true "ID поста"
// @Success 200 {object} models.DeletePostResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [delete]
func (h *Handler) DeletePostHandler(c *gin.Context) {
	userID, exists := c.Get(UserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User ID not found"})
		return
	}

	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Post ID is required"})
		return
	}

	req := &postpb.DeletePostRequest{
		Id:        postID,
		CreatorId: userID.(string),
	}

	resp, err := h.PostServiceClient.DeletePost(c, req)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, models.DeletePostResponse{Success: resp.Success})
}

// ListPostsHandler обрабатывает запросы на получение списка постов
// @Summary Получение списка постов
// @Description Возвращает пагинированный список постов
// @Tags Posts
// @Produce json
// @Security sessionAuth
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param page_size query int false "Размер страницы (по умолчанию 20)"
// @Success 200 {object} models.PostListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts [get]
func (h *Handler) ListPostsHandler(c *gin.Context) {
	userID, exists := c.Get(UserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User ID not found"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	req := &postpb.ListPostsRequest{
		Page:        int32(page),
		PageSize:    int32(pageSize),
		RequesterId: userID.(string),
	}

	resp, err := h.PostServiceClient.ListPosts(c, req)
	if err != nil {
		log.Printf("Error listing posts: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list posts"})
		return
	}

	posts := resp.GetPosts()
	postResponses := make([]models.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = convertProtoToResponse(post)
	}

	c.JSON(http.StatusOK, models.PostListResponse{
		Posts:      postResponses,
		TotalCount: int(resp.TotalCount),
	})
}

func convertProtoToResponse(post *postpb.Post) models.PostResponse {
	postIDInt, _ := strconv.ParseUint(post.Id, 10, 32)
	creatorIDInt, _ := strconv.ParseUint(post.CreatorId, 10, 32)

	return models.PostResponse{
		ID:          uint(postIDInt),
		Title:       post.Title,
		Description: post.Description,
		CreatorID:   uint(creatorIDInt),
		CreatedAt:   post.CreatedAt.AsTime(),
		UpdatedAt:   post.UpdatedAt.AsTime(),
		IsPrivate:   post.IsPrivate,
		Tags:        post.Tags,
	}
}
