package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lapkomo2018/goTwitterServices/internal/post/model"
	"github.com/lapkomo2018/goTwitterServices/pkg/constant"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	imageMaxSize = 50 * 1024 * 1024  // 50MB
	videoMaxSize = 100 * 1024 * 1024 // 100MB
)

func (h *Handler) initPostApi(api *echo.Group) {
	posts := api.Group("/posts")
	{
		posts.POST("", h.createPost, h.authenticationMiddleware, h.banMiddleware)

		postID := posts.Group(fmt.Sprintf("/:%s", postIDParam), h.setPostByIDMiddleware)
		{
			postID.GET("", h.getPost)
			postID.PATCH("", h.updatePost, h.authenticationMiddleware, h.banMiddleware)
			postID.DELETE("", h.deletePost, h.authenticationMiddleware, h.banMiddleware)
		}
	}
}

func (h *Handler) createPost(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*user)
	if !ok {
		return echo.ErrUnauthorized
	}

	// Parse the form data
	if err := c.Request().ParseMultipartForm(32 << 20); err != nil {
		return nil
	}

	req := &struct {
		Title   string                  `json:"title" form:"title" validate:"required,min=1,max=255"`
		Content string                  `json:"content" form:"content" validate:"required,min=1,max=65535"`
		Tags    []string                `json:"tags" form:"tags" validate:"omitempty,dive,min=1"`
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
	requester, ok := c.Get(requesterLocals).(*user)
	if !ok {
		return echo.ErrUnauthorized
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

	// Process images
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

	// Process videos
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

	post.UpdatedBy = requester.ID

	if err := h.ps.Update(post); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, post)
}

func (h *Handler) deletePost(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*user)
	if !ok {
		return echo.ErrUnauthorized
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
