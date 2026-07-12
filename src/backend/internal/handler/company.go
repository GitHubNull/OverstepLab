package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/oversteplab/oversteplab/internal/common"
	"github.com/oversteplab/oversteplab/internal/middleware"
	"github.com/oversteplab/oversteplab/internal/service"
)

type CompanyHandler struct {
	companyService *service.CompanyService
}

func NewCompanyHandler(svc *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: svc}
}

func (h *CompanyHandler) ListMembers(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	members, err := h.companyService.ListMembers(user)
	if err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.InternalError(c, err.Error())
		}
		return
	}
	common.Success(c, members)
}

func (h *CompanyHandler) AddMember(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	var input service.AddMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}
	member, err := h.companyService.AddMember(user, &input)
	if err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.InternalError(c, err.Error())
		}
		return
	}
	common.Success(c, member)
}

func (h *CompanyHandler) UpdateMember(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	var input service.UpdateMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}
	if err := h.companyService.UpdateMember(user, input.ID, &input); err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else if err == service.ErrMemberNotFound {
			common.NotFound(c, "Member not found")
		} else {
			common.InternalError(c, err.Error())
		}
		return
	}
	common.SuccessMessage(c, "Member updated")
}

func (h *CompanyHandler) DeleteMember(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	var input struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}
	if err := h.companyService.DeleteMember(user, input.ID); err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else if err == service.ErrMemberNotFound {
			common.NotFound(c, "Member not found")
		} else {
			common.InternalError(c, err.Error())
		}
		return
	}
	common.SuccessMessage(c, "Member deleted")
}

func (h *CompanyHandler) ChangeRole(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	var input struct {
		ID   uint   `json:"id"`
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Role == "" {
		common.BadRequest(c, "Invalid request")
		return
	}
	if err := h.companyService.ChangeRole(user, input.ID, input.Role); err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else if err == service.ErrMemberNotFound {
			common.NotFound(c, "Member not found")
		} else {
			common.InternalError(c, err.Error())
		}
		return
	}
	common.SuccessMessage(c, "Role updated")
}
