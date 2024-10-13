package v1

import (
	"fmt"
	"github.com/SocialNetworkY/Backend/internal/post/model"
	"github.com/SocialNetworkY/Backend/pkg/constant"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) initCommentsApi(api *echo.Group) {
	comments := api.Group("/comments")
	{
		commentID := comments.Group(fmt.Sprintf("/:%s", commentIDParam), h.setCommentByIDMiddleware)
		{
			commentID.PUT("", h.changeComment, h.authenticationMiddleware)
			commentID.DELETE("", h.deleteComment, h.authenticationMiddleware)
		}

	}
}

func (h *Handler) changeComment(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get requester from context")
	}

	comment, ok := c.Get(commentLocals).(*model.Comment)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get comment from context")
	}

	if requester.ID != comment.UserID && requester.Role == constant.RoleUser {
		return echo.NewHTTPError(http.StatusForbidden, "you are not allowed to edit this comment")
	}

	var req = struct {
		Content string `json:"content" validate:"required,min=1,max=255"`
	}{}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.cs.EditComment(comment.ID, requester.ID, req.Content); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	comment, err := h.cs.Find(comment.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, comment)
}

func (h *Handler) deleteComment(c echo.Context) error {
	requester, ok := c.Get(requesterLocals).(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get requester from context")
	}

	comment, ok := c.Get(commentLocals).(*model.Comment)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get comment from context")
	}

	if requester.ID != comment.UserID && requester.Role == constant.RoleUser {
		return echo.NewHTTPError(http.StatusForbidden, "you are not allowed to delete this comment")
	}

	if err := h.cs.DeleteComment(comment.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
