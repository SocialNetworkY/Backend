package v1

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/SocialNetworkY/Backend/internal/post/model"
	"github.com/SocialNetworkY/Backend/pkg/constant"
	"github.com/labstack/echo/v4"
)

const (
	imageMaxSize = 50 * 1024 * 1024  // 50MB
	videoMaxSize = 100 * 1024 * 1024 // 100MB
)

func (h *Handler) initPostApi(api *echo.Group) {
	posts := api.Group("/posts")
	{
		posts.GET("", h.getPosts)
		posts.GET("/search", h.searchPosts)
		posts.POST("", h.createPost, h.authenticationMiddleware, h.banMiddleware)
		posts.GET(fmt.Sprintf("/users/:%s", userIDParam), h.getPostsByUserID, h.setUserByIDMiddleware)

		postID := posts.Group(fmt.Sprintf("/:%s", postIDParam), h.setPostByIDMiddleware)
		{
			postID.GET("", h.getPost)
			postID.PATCH("", h.updatePost, h.authenticationMiddleware, h.banMiddleware)
			postID.DELETE("", h.deletePost, h.authenticationMiddleware, h.banMiddleware)
			postID.GET("/like", h.likePost, h.authenticationMiddleware, h.banMiddleware)
			postID.GET("/unlike", h.unlikePost, h.authenticationMiddleware, h.banMiddleware)

			// Comments
			comments := postID.Group("/comments")
			{
				comments.GET("", h.getPostComments)
				comments.POST("", h.commentPost, h.authenticationMiddleware, h.banMiddleware)
			}
		}
	}
}

func (h *Handler) createPost(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get requester")
	}

	req := &struct {
		Title   string                  `form:"title" validate:"required,min=1,max=255"`
		Content string                  `form:"content" validate:"required,min=1,max=65535"`
		Tags    []string                `form:"tags" validate:"omitempty,dive,min=1"`
		Images  []*multipart.FileHeader `form:"images"`
		Videos  []*multipart.FileHeader `form:"videos"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Process images
	var imageUrls []string
	for _, imageHeader := range req.Images {
		if !strings.HasPrefix(imageHeader.Header.Get("Content-Type"), "image/") {
			return echo.NewHTTPError(http.StatusBadRequest, "file is not an image")
		}

		if imageHeader.Size > imageMaxSize {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("image file size is too large, max size is %d", imageMaxSize))
		}

		image, err := imageHeader.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		imageUrl, err := h.fs.UploadFile(image, fmt.Sprintf("%d_%d_%s", requester.ID, time.Now().Unix(), imageHeader.Filename))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		imageUrls = append(imageUrls, imageUrl)
	}

	// Process videos
	var videoUrls []string
	for _, videoHeader := range req.Videos {
		if !strings.HasPrefix(videoHeader.Header.Get("Content-Type"), "video/") {
			return echo.NewHTTPError(http.StatusBadRequest, "file is not a video")
		}

		if videoHeader.Size > videoMaxSize {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("video file size is too large, max size is %d", videoMaxSize))
		}

		video, err := videoHeader.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		videoUrl, err := h.fs.UploadFile(video, fmt.Sprintf("%d_%d_%s", requester.ID, time.Now().Unix(), videoHeader.Filename))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		videoUrls = append(videoUrls, videoUrl)
	}

	tags := make([]*model.Tag, len(req.Tags))
	for i, tag := range req.Tags {
		tags[i] = &model.Tag{
			Name: tag,
		}
	}

	post := &model.Post{
		UserID:    requester.ID,
		Title:     req.Title,
		Content:   req.Content,
		Tags:      tags,
		ImageURLs: imageUrls,
		VideoURLs: videoUrls,
		PostedAt:  time.Now(),
	}

	if err := h.ps.Create(post); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, post)
}

func (h *Handler) getPost(c echo.Context) error {
	post, ok := c.Get(postLocals).(*model.Post)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get post")
	}

	return c.JSON(http.StatusOK, post)
}

func (h *Handler) updatePost(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get requester")
	}

	post, ok := c.Get(postLocals).(*model.Post)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get post")
	}

	if requester.ID != post.UserID && requester.Role < constant.RoleAdminLvl1 {
		return echo.NewHTTPError(http.StatusForbidden, "you are not allowed to update this post")
	}

	req := &struct {
		Title     string                  `form:"title" validate:"omitempty,min=1,max=255"`
		Content   string                  `form:"content" validate:"omitempty,min=1,max=65535"`
		Tags      []string                `form:"tags" validate:"omitempty,dive,min=1"`
		ImageUrls []string                `form:"image_urls" validate:"omitempty,dive,url"`
		VideoUrls []string                `form:"video_urls" validate:"omitempty,dive,url"`
		Images    []*multipart.FileHeader `form:"images"`
		Videos    []*multipart.FileHeader `form:"videos"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if req.Title != "" {
		post.Title = req.Title
	}

	if req.Content != "" {
		post.Content = req.Content
	}

	if len(req.Tags) > 0 {
		tags := make([]*model.Tag, len(req.Tags))
		for i, tag := range req.Tags {
			tags[i] = &model.Tag{
				Name: tag,
			}
		}
		post.Tags = tags
	}

	if len(req.ImageUrls) > 0 {
		post.ImageURLs = req.ImageUrls
	}

	for _, imageHeader := range req.Images {
		if imageHeader.Size > imageMaxSize {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("image file size is too large, max size is %d", imageMaxSize))
		}

		image, err := imageHeader.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		imageUrl, err := h.fs.UploadFile(image, fmt.Sprintf("%d_%d_%s", requester.ID, time.Now().Unix(), imageHeader.Filename))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		post.ImageURLs = append(post.ImageURLs, imageUrl)
	}

	if len(req.VideoUrls) > 0 {
		post.VideoURLs = req.VideoUrls
	}

	for _, videoHeader := range req.Videos {
		if videoHeader.Size > videoMaxSize {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("video file size is too large, max size is %d", videoMaxSize))
		}

		video, err := videoHeader.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		videoUrl, err := h.fs.UploadFile(video, fmt.Sprintf("%d_%d_%s", requester.ID, time.Now().Unix(), videoHeader.Filename))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		post.VideoURLs = append(post.VideoURLs, videoUrl)
	}

	post.EditedBy = requester.ID

	if err := h.ps.Update(post); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, post)
}

func (h *Handler) deletePost(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get requester")
	}

	post, ok := c.Get(postLocals).(*model.Post)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get post")
	}

	if requester.ID != post.UserID && requester.Role < constant.RoleAdminLvl1 {
		return echo.NewHTTPError(http.StatusForbidden, "you are not allowed to delete this post")
	}

	if err := h.ps.Delete(post.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) getPosts(c echo.Context) error {
	skip, limit := skipLimitQuery(c)

	posts, err := h.ps.FindSome(skip, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, posts)
}

func (h *Handler) getPostsByUserID(c echo.Context) error {
	user, ok := c.Get(userLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

	skip, limit := skipLimitQuery(c)

	posts, err := h.ps.FindByUser(user.ID, skip, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, posts)
}

func (h *Handler) likePost(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get requester")
	}

	post, ok := c.Get(postLocals).(*model.Post)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get post")
	}

	if err := h.ls.LikePost(post.ID, requester.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) unlikePost(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get requester")
	}

	post, ok := c.Get(postLocals).(*model.Post)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get post")
	}

	if err := h.ls.UnlikePost(post.ID, requester.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) getPostComments(c echo.Context) error {
	post, ok := c.Get(postLocals).(*model.Post)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get post")
	}

	skip, limit := skipLimitQuery(c)

	comments, err := h.cs.FindByPost(post.ID, skip, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, comments)
}

func (h *Handler) commentPost(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get requester")
	}

	if requester.Banned {
		return echo.NewHTTPError(http.StatusForbidden, "you are banned")
	}

	post, ok := c.Get(postLocals).(*model.Post)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get post")
	}

	req := &struct {
		Content string `json:"content" validate:"required,min=1,max=255"`
	}{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.cs.CommentPost(post.ID, requester.ID, req.Content); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) searchPosts(c echo.Context) error {
	skip, limit := skipLimitQuery(c)

	query := c.QueryParam(queryQuery)
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "query is required")
	}

	posts, err := h.ps.Search(query, skip, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, posts)
}
