package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oversteplab/oversteplab/internal/common"
	"github.com/oversteplab/oversteplab/internal/middleware"
	"github.com/oversteplab/oversteplab/internal/service"
)

type VPSHandler struct {
	vpsService *service.VPSService
}

func NewVPSHandler(svc *service.VPSService) *VPSHandler {
	return &VPSHandler{vpsService: svc}
}

func (h *VPSHandler) List(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	vpsList, err := h.vpsService.List(user)
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, vpsList)
}

func (h *VPSHandler) Create(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	var input service.CreateVPSInput
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}
	vps, err := h.vpsService.Create(user, &input)
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, vps)
}

func (h *VPSHandler) GetByID(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	vps, err := h.vpsService.GetDetail(user, uint(id))
	if err != nil {
		common.NotFound(c, "VPS not found")
		return
	}
	common.Success(c, vps)
}

func (h *VPSHandler) Start(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.vpsService.StartVPS(user, uint(id)); err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.NotFound(c, err.Error())
		}
		return
	}
	common.SuccessMessage(c, "VPS started")
}

func (h *VPSHandler) Stop(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.vpsService.StopVPS(user, uint(id)); err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.NotFound(c, err.Error())
		}
		return
	}
	common.SuccessMessage(c, "VPS stopped")
}

func (h *VPSHandler) Restart(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.vpsService.RestartVPS(user, uint(id)); err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.NotFound(c, err.Error())
		}
		return
	}
	common.SuccessMessage(c, "VPS restarted")
}

func (h *VPSHandler) Reinstall(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var input struct {
		OSImage string `json:"os_image"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}
	if err := h.vpsService.ReinstallVPS(user, uint(id), input.OSImage); err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.NotFound(c, err.Error())
		}
		return
	}
	common.SuccessMessage(c, "VPS reinstalled")
}

func (h *VPSHandler) Delete(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.vpsService.DeleteVPS(user, uint(id)); err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.NotFound(c, err.Error())
		}
		return
	}
	common.SuccessMessage(c, "VPS deleted")
}

func (h *VPSHandler) Console(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	console, err := h.vpsService.GetConsole(user, uint(id))
	if err != nil {
		common.NotFound(c, err.Error())
		return
	}
	common.Success(c, console)
}
