package handlers

import (
	"log"
	"net/http"
	"social-network/common/models"
	"strconv"

	postpb "social-network/common/proto/post"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func handleGRPCError(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Unknown internal error"})
		return
	}

	httpStatus := grpcCodeToHTTPStatus(st.Code())
	c.JSON(httpStatus, models.ErrorResponse{Error: st.Message()})
}

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
		CreatorId:   userID.(uint64),
		IsPrivate:   reqBody.IsPrivate,
		Tags:        reqBody.Tags,
	}

	resp, err := h.PostServiceClient.CreatePost(c, req)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		handleGRPCError(c, err)
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

	postIDInt, _ := strconv.ParseUint(postID, 10, 64)

	req := &postpb.GetPostRequest{
		Id:          postIDInt,
		RequesterId: userID.(uint64),
	}

	resp, err := h.PostServiceClient.GetPost(c, req)
	if err != nil {
		log.Printf("Error getting post: %v", err)
		handleGRPCError(c, err)
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

	postIDInt, _ := strconv.ParseUint(postID, 10, 64)

	req := &postpb.UpdatePostRequest{
		Id:          postIDInt,
		Title:       reqBody.Title,
		Description: reqBody.Description,
		CreatorId:   userID.(uint64),
		IsPrivate:   reqBody.IsPrivate,
		Tags:        reqBody.Tags,
	}

	resp, err := h.PostServiceClient.UpdatePost(c, req)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		handleGRPCError(c, err)
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

	postIDInt, _ := strconv.ParseUint(postID, 10, 64)

	req := &postpb.DeletePostRequest{
		Id:        postIDInt,
		CreatorId: userID.(uint64),
	}

	resp, err := h.PostServiceClient.DeletePost(c, req)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		handleGRPCError(c, err)
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
		RequesterId: userID.(uint64),
	}

	resp, err := h.PostServiceClient.ListPosts(c, req)
	if err != nil {
		log.Printf("Error listing post: %v", err)
		handleGRPCError(c, err)
		return
	}

	posts := resp.GetPosts()
	postResponses := make([]models.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = convertProtoToResponse(post)
	}

	c.JSON(http.StatusOK, models.PostListResponse{
		Posts:      postResponses,
		TotalCount: resp.TotalCount,
	})
}

func convertProtoToResponse(post *postpb.Post) models.PostResponse {
	return models.PostResponse{
		ID:          post.Id,
		Title:       post.Title,
		Description: post.Description,
		CreatorID:   post.CreatorId,
		CreatedAt:   post.CreatedAt.AsTime(),
		UpdatedAt:   post.UpdatedAt.AsTime(),
		IsPrivate:   post.IsPrivate,
		Tags:        post.Tags,
	}
}

// ViewPostHandler обрабатывает запросы на просмотр поста
// @Summary Просмотр поста
// @Description Регистрирует просмотр поста текущим пользователем и отправляет событие в Kafka
// @Tags Posts
// @Produce json
// @Security sessionAuth
// @Param id path string true "ID поста"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/view [post]
func (h *Handler) ViewPostHandler(c *gin.Context) {
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

	postIDInt, _ := strconv.ParseUint(postID, 10, 64)

	req := &postpb.ViewPostRequest{
		Id:     postIDInt,
		UserId: userID.(uint64),
	}

	resp, err := h.PostServiceClient.ViewPost(c, req)
	if err != nil {
		log.Printf("Error viewing post: %v", err)
		handleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: resp.Success})
}

// LikePostHandler обрабатывает запросы на лайк поста
// @Summary Лайк поста
// @Description Ставит или убирает лайк поста текущим пользователем и отправляет событие в Kafka
// @Tags Posts
// @Produce json
// @Security sessionAuth
// @Param id path string true "ID поста"
// @Success 200 {object} models.LikePostResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/like [post]
func (h *Handler) LikePostHandler(c *gin.Context) {
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

	postIDInt, _ := strconv.ParseUint(postID, 10, 64)

	req := &postpb.LikePostRequest{
		Id:     postIDInt,
		UserId: userID.(uint64),
	}

	resp, err := h.PostServiceClient.LikePost(c, req)
	if err != nil {
		log.Printf("Error liking post: %v", err)
		handleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.LikePostResponse{
		Success:    resp.Success,
		TotalLikes: int(resp.TotalLikes),
	})
}

// AddCommentHandler обрабатывает запросы на добавление комментария к посту
// @Summary Добавление комментария к посту
// @Description Добавляет новый комментарий к посту от имени текущего пользователя и отправляет событие в Kafka
// @Tags Comments
// @Accept json
// @Produce json
// @Security sessionAuth
// @Param id path string true "ID поста"
// @Param body body models.AddCommentRequest true "Данные комментария"
// @Success 201 {object} models.CommentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/comments [post]
func (h *Handler) AddCommentHandler(c *gin.Context) {
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

	var reqBody models.AddCommentRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	postIDint, _ := strconv.ParseUint(postID, 10, 64)

	req := &postpb.AddCommentRequest{
		PostId:  postIDint,
		UserId:  userID.(uint64),
		Content: reqBody.Content,
	}

	resp, err := h.PostServiceClient.AddComment(c, req)
	if err != nil {
		log.Printf("Error adding comment: %v", err)
		handleGRPCError(c, err)
		return
	}

	comment := models.Comment{
		ID:        resp.Comment.Id,
		PostID:    resp.Comment.PostId,
		UserID:    resp.Comment.UserId,
		Content:   resp.Comment.Content,
		CreatedAt: resp.Comment.CreatedAt.AsTime(),
	}

	c.JSON(http.StatusCreated, models.CommentResponse{
		Comment: comment,
	})
}

// GetCommentsHandler обрабатывает запросы на получение комментариев к посту
// @Summary Получение комментариев к посту
// @Description Возвращает список комментариев к посту с пагинацией
// @Tags Comments
// @Produce json
// @Security sessionAuth
// @Param id path string true "ID поста"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param page_size query int false "Размер страницы (по умолчанию 20)"
// @Success 200 {object} models.CommentListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/comments [get]
func (h *Handler) GetCommentsHandler(c *gin.Context) {
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

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseUint(pageSizeStr, 10, 64)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	postIDInt, _ := strconv.ParseUint(postID, 10, 64)

	req := &postpb.GetCommentsRequest{
		PostId:      postIDInt,
		Page:        page,
		PageSize:    pageSize,
		RequesterId: userID.(uint64),
	}

	resp, err := h.PostServiceClient.GetComments(c, req)
	if err != nil {
		log.Printf("Error getting comments: %v", err)
		handleGRPCError(c, err)
		return
	}

	comments := make([]models.Comment, len(resp.Comments))
	for i, comment := range resp.Comments {
		comments[i] = models.Comment{
			ID:        comment.Id,
			PostID:    comment.PostId,
			UserID:    comment.UserId,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.AsTime(),
		}
	}

	c.JSON(http.StatusOK, models.CommentListResponse{
		Comments:   comments,
		TotalCount: resp.TotalCount,
	})
}
