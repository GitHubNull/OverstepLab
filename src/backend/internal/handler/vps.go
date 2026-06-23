package handler

import (
	"fmt"
	"net/http"
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

func handleVPSError(c *gin.Context, err error) {
	switch err {
	case service.ErrUnauthorized:
		common.Forbidden(c, err.Error())
	case service.ErrVPSNotFound:
		common.NotFound(c, err.Error())
	default:
		common.InternalError(c, err.Error())
	}
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
		handleVPSError(c, err)
		return
	}
	common.SuccessMessage(c, "VPS started")
}

func (h *VPSHandler) Stop(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.vpsService.StopVPS(user, uint(id)); err != nil {
		handleVPSError(c, err)
		return
	}
	common.SuccessMessage(c, "VPS stopped")
}

func (h *VPSHandler) Restart(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.vpsService.RestartVPS(user, uint(id)); err != nil {
		handleVPSError(c, err)
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
		handleVPSError(c, err)
		return
	}
	common.SuccessMessage(c, "VPS reinstalled")
}

func (h *VPSHandler) Delete(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.vpsService.DeleteVPS(user, uint(id)); err != nil {
		handleVPSError(c, err)
		return
	}
	common.SuccessMessage(c, "VPS deleted")
}

func (h *VPSHandler) Console(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	console, err := h.vpsService.GetConsole(user, uint(id))
	if err != nil {
		handleVPSError(c, err)
		return
	}
	common.Success(c, console)
}

func (h *VPSHandler) ConsoleView(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		common.Unauthorized(c, "Missing console token")
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	data, err := h.vpsService.GetConsoleView(token, uint(id))
	if err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.NotFound(c, err.Error())
		}
		return
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><title>Console - VPS %s</title></head>
<body>
  <h1>VPS 控制台</h1>
  <p><strong>状态:</strong> %s</p>
  <p><strong>IP:</strong> %s</p>
  <p><strong>WebSocket:</strong> %s</p>
  <p><strong>Token:</strong> %s</p>
</body>
</html>`, data["vps_id"], data["status"], data["ip_address"], data["url"], data["token"])

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
