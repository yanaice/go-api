package handler

import (
	"github.com/gin-gonic/gin"
	"go-starter-project/internal/app/model"
	"go-starter-project/pkg/auth"
	"go-starter-project/pkg/derror"
	commonhdl "go-starter-project/pkg/handler"
	"net/http"
)

func (h *Handler) registerTagSvc(r *gin.RouterGroup) {
	for _, route := range h.TagRoutes() {
		r.Handle(route.Method, route.Pattern, route.Handler)
	}
}

func (h *Handler) TagRoutes() []Route {
	return []Route{
		{
			Name:        "CreateTag",
			Description: "",
			Method:      http.MethodPost,
			Pattern:     "/",
			Handler:     h.CreateTag,
		},
		{
			Name:        "UpdateTag",
			Description: "",
			Method:      http.MethodPut,
			Pattern:     "/",
			Handler:     h.UpdateTag,
		},
		{
			Name:        "DeleteTag",
			Description: "",
			Method:      http.MethodDelete,
			Pattern:     "/",
			Handler:     h.DeleteTag,
		},
		{
			Name:        "ReadTag",
			Description: "",
			Method:      http.MethodGet,
			Pattern:     "/:id",
			Handler:     h.ReadTag,
		},
		{
			Name:        "ReadTags",
			Description: "",
			Method:      http.MethodGet,
			Pattern:     "/",
			Handler:     h.ReadTags,
		},
	}
}

func (h *Handler) CreateTag(c *gin.Context) {
	doer, _ := c.Get(auth.PetStoreUser)

	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		commonhdl.HandlerReturnError(c, derror.E(err).SetCode(derror.ErrCodeInputError).
			SetDebug("Create tag: cannot read json").Log())
		return
	}
	if err := h.TagSvc.CreateTag(doer, tag); err != nil {
		commonhdl.HandlerReturnError(c, err)
		return
	}
	commonhdl.HandlerReturnData(c, model.Response{Status: model.ResponseStatusSuccess})
}

func (h *Handler) UpdateTag(c *gin.Context) {
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		commonhdl.HandlerReturnError(c, derror.E(err).SetCode(derror.ErrCodeInputError).
			SetDebug("Update tag: cannot read json").Log())
		return
	}
	if err := h.TagSvc.UpdateTag(tag.ID, tag.Name); err != nil {
		c.JSON(http.StatusInternalServerError, "Internal error")
		return
	}
	commonhdl.HandlerReturnData(c, model.Response{Status: model.ResponseStatusSuccess})
}

func (h *Handler) DeleteTag(c *gin.Context) {
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		commonhdl.HandlerReturnError(c, derror.E(err).SetCode(derror.ErrCodeInputError).
			SetDebug("Delete tag: cannot read json").Log())
		return
	}
	if err := h.TagSvc.DeleteTag(tag.ID); err != nil {
		commonhdl.HandlerReturnError(c, err)
		return
	}
	commonhdl.HandlerReturnData(c, model.Response{Status: model.ResponseStatusSuccess})
}

func (h *Handler) ReadTag(c *gin.Context) {
	tagID := c.Param("id")
	res, err := h.TagSvc.ReadTag(tagID)
	if err != nil {
		commonhdl.HandlerReturnError(c, err)
		return
	}
	commonhdl.HandlerReturnData(c, res)
}

func (h *Handler) ReadTags(c *gin.Context) {
	res, err := h.TagSvc.ReadTags()
	if err != nil {
		commonhdl.HandlerReturnError(c, err)
		return
	}
	commonhdl.HandlerReturnData(c, res)
}
