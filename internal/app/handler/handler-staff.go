package handler

import (
	"github.com/gin-gonic/gin"
	"go-starter-project/internal/app/model"
	"go-starter-project/pkg/auth"
	"go-starter-project/pkg/derror"
	commonhdl "go-starter-project/pkg/handler"
	"net/http"
)

func (h *Handler) registerStaffSvc(r *gin.RouterGroup) {
	for _, route := range h.staffSvcRoutes() {
		r.Handle(route.Method, route.Pattern, route.Handler)
	}
}

func (h *Handler) staffSvcRoutes() []Route {
	return []Route{
		{
			Name:        "LoginStaffUser",
			Description: "Login staff user",
			Method:      http.MethodPost,
			Pattern:     "/login",
			Handler:     h.loginStaffUser,
		},
		{
			Name:        "CreateUser",
			Description: "Create staff user",
			Method:      http.MethodPost,
			Pattern:     "/create",
			Handler:     h.CreateUser,
		},
	}
}

func (h *Handler) loginStaffUser(c *gin.Context) {
	var login model.StaffLoginRequest

	if err := c.ShouldBindJSON(&login); err != nil {
		commonhdl.HandlerReturnError(c, derror.E(err).SetCode(derror.ErrCodeInputError).
			SetDebug("login failed: bind body error").Log())
		return
	}

	tokens, err := h.StaffSvc.LoginUser(login.Username, login.Password, c.ClientIP())
	if err != nil {
		commonhdl.HandlerReturnError(c, err)
		return
	}

	commonhdl.HandlerReturnData(c, tokens)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var login model.StaffLoginRequest

	doer, _ := c.Get(auth.PetStoreUser)

	if err := c.ShouldBindJSON(&login); err != nil {
		commonhdl.HandlerReturnError(c, derror.E(err).SetCode(derror.ErrCodeInputError).
			SetDebug("login failed: bind body error").Log())
		return
	}

	err := h.StaffSvc.CreateUser(doer, login.Username, login.Password, c.ClientIP())
	if err != nil {
		commonhdl.HandlerReturnError(c, err)
		return
	}

	commonhdl.HandlerReturnData(c, model.Response{Status: model.ResponseStatusSuccess})
}
