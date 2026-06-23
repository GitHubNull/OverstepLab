package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oversteplab/oversteplab/database"
	"github.com/oversteplab/oversteplab/database/seed"
	"github.com/oversteplab/oversteplab/internal/common"
	"github.com/oversteplab/oversteplab/internal/middleware"
	"github.com/oversteplab/oversteplab/internal/repository"
	"github.com/oversteplab/oversteplab/internal/service"
	"github.com/oversteplab/oversteplab/internal/vuln"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: svc}
}

func (h *OrderHandler) List(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	orders, err := h.orderService.List(user)
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, orders)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	order, err := h.orderService.GetDetail(user, uint(id))
	if err != nil {
		common.NotFound(c, "Order not found")
		return
	}
	common.Success(c, order)
}

type BillHandler struct {
	billService *service.BillService
}

func NewBillHandler(svc *service.BillService) *BillHandler {
	return &BillHandler{billService: svc}
}

func (h *BillHandler) List(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	bills, err := h.billService.List(user)
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, bills)
}

func (h *BillHandler) Recharge(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	var input struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Amount <= 0 {
		common.BadRequest(c, "Invalid amount")
		return
	}
	if err := h.billService.Recharge(user, input.Amount); err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.SuccessMessage(c, "Recharge successful")
}

func (h *BillHandler) Export(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	bills, err := h.billService.Export(user)
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}

	// Generate CSV
	csv := "ID,Type,Amount,Balance After,Description,Date\n"
	for _, b := range bills {
		csv += strconv.FormatUint(uint64(b.ID), 10) + "," +
			b.Type + "," +
			strconv.FormatFloat(b.Amount, 'f', 2, 64) + "," +
			strconv.FormatFloat(b.BalanceAfter, 'f', 2, 64) + "," +
			b.Description + "," +
			b.CreatedAt.Format("2006-01-02 15:04:05") + "\n"
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=bills.csv")
	c.Data(http.StatusOK, "text/csv", []byte(csv))
}

type TicketHandler struct {
	ticketService *service.TicketService
}

func NewTicketHandler(svc *service.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: svc}
}

func (h *TicketHandler) List(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	tickets, err := h.ticketService.List(user)
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, tickets)
}

func (h *TicketHandler) Create(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	var input service.CreateTicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}
	ticket, err := h.ticketService.Create(user, &input)
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, ticket)
}

func (h *TicketHandler) GetByID(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ticket, err := h.ticketService.GetDetail(user, uint(id))
	if err != nil {
		common.NotFound(c, "Ticket not found")
		return
	}
	replies, _ := h.ticketService.GetReplies(uint(id))
	common.Success(c, gin.H{
		"ticket": ticket,
		"replies": replies,
	})
}

func (h *TicketHandler) Reply(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var input struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Content == "" {
		common.BadRequest(c, "Invalid request")
		return
	}
	reply, err := h.ticketService.Reply(user, uint(id), input.Content)
	if err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.InternalError(c, err.Error())
		}
		return
	}
	common.Success(c, reply)
}

func (h *TicketHandler) Close(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.ticketService.Close(user, uint(id)); err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.NotFound(c, err.Error())
		}
		return
	}
	common.SuccessMessage(c, "Ticket closed")
}

type APIKeyHandler struct {
	apiKeyService *service.APIKeyService
}

func NewAPIKeyHandler(svc *service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{apiKeyService: svc}
}

func (h *APIKeyHandler) List(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	keys, err := h.apiKeyService.List(user)
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, keys)
}

func (h *APIKeyHandler) Create(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	var input struct {
		Name        string `json:"name"`
		Permissions string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}
	key, err := h.apiKeyService.Create(user, input.Name, input.Permissions)
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	// Return raw key in response (only time it's ever shown)
	common.Success(c, gin.H{
		"id":          key.ID,
		"user_id":     key.UserID,
		"name":        key.Name,
		"key_value":   key.KeyValue, // raw key, only shown on creation
		"key_prefix":  key.KeyPrefix,
		"permissions": key.Permissions,
		"status":      key.Status,
		"last_used_at": key.LastUsedAt,
		"expire_at":   key.ExpireAt,
		"created_at":  key.CreatedAt,
	})
}

func (h *APIKeyHandler) Delete(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.apiKeyService.Delete(user, uint(id)); err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.NotFound(c, err.Error())
		}
		return
	}
	common.SuccessMessage(c, "API key deleted")
}

type AuditLogHandler struct {
	auditLogService *service.AuditLogService
}

func NewAuditLogHandler(svc *service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{auditLogService: svc}
}

func (h *AuditLogHandler) List(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	logs, err := h.auditLogService.GetLogs(user)
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, logs)
}

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(svc *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: svc}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.adminService.ListUsers()
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, users)
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var input struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}
	if err := h.adminService.UpdateUserStatus(uint(id), input.Status); err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.SuccessMessage(c, "User status updated")
}

func (h *AdminHandler) ResetUserPassword(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var input struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || len(input.Password) < 6 {
		common.BadRequest(c, "Invalid password (minimum 6 characters)")
		return
	}
	if err := h.adminService.ResetUserPassword(uint(id), input.Password); err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.SuccessMessage(c, "Password reset successful")
}

func (h *AdminHandler) ListCompanies(c *gin.Context) {
	companies, err := h.adminService.ListCompanies()
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, companies)
}

func (h *AdminHandler) ListAllVPS(c *gin.Context) {
	vpsList, err := h.adminService.ListAllVPS()
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, vpsList)
}

func (h *AdminHandler) ListAllLogs(c *gin.Context) {
	logs, err := h.adminService.ListAllLogs()
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, logs)
}

type ChallengeHandler struct{
	dbPath string
}

func NewChallengeHandler(dbPath string) *ChallengeHandler {
	return &ChallengeHandler{dbPath: dbPath}
}

func (h *ChallengeHandler) List(c *gin.Context) {
	challenges := make([]map[string]interface{}, 0)
	for _, ch := range vuln.Challenges {
		challenges = append(challenges, map[string]interface{}{
			"id":          ch.ID,
			"title":       ch.Title,
			"category":    ch.Category,
			"difficulty":  ch.Difficulty,
			"description": ch.Description,
			"completed":   vuln.IsChallengeCompleted(ch.ID),
		})
	}
	common.Success(c, challenges)
}

func (h *ChallengeHandler) Detail(c *gin.Context) {
	id := c.Param("id")
	for _, ch := range vuln.Challenges {
		if ch.ID == id {
			// Filter out sensitive fields (hints and writeup) from detail response
			common.Success(c, map[string]interface{}{
				"id":          ch.ID,
				"title":       ch.Title,
				"category":    ch.Category,
				"difficulty":  ch.Difficulty,
				"description": ch.Description,
				"endpoint":    ch.Endpoint,
				"method":      ch.Method,
				"completed":   vuln.IsChallengeCompleted(ch.ID),
			})
			return
		}
	}
	common.NotFound(c, "Challenge not found")
}

func (h *ChallengeHandler) GetHint(c *gin.Context) {
	id := c.Param("id")
	level, _ := strconv.ParseUint(c.Param("level"), 10, 64)
	for _, ch := range vuln.Challenges {
		if ch.ID == id {
			if level == 0 || int(level) > len(ch.Hints) {
				common.BadRequest(c, "Invalid hint level")
				return
			}
			common.Success(c, gin.H{
				"hint": ch.Hints[level-1],
				"level": level,
			})
			return
		}
	}
	common.NotFound(c, "Challenge not found")
}

func (h *ChallengeHandler) MarkComplete(c *gin.Context) {
	id := c.Param("id")
	vuln.MarkChallengeCompleted(id)
	common.SuccessMessage(c, "Challenge marked as completed")
}

func (h *ChallengeHandler) GetSecurityMode(c *gin.Context) {
	common.Success(c, gin.H{
		"mode":        vuln.GetModeString(),
		"description": "Current security mode of the application",
	})
}

func (h *ChallengeHandler) SetSecurityMode(c *gin.Context) {
	var input struct {
		Mode string `json:"mode"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}
	vuln.SetSecureMode(input.Mode == "secure")
	common.SuccessMessage(c, "Security mode updated to " + vuln.GetModeString())
}

func (h *ChallengeHandler) ResetDatabase(c *gin.Context) {
	vuln.ResetProgress()
	if err := database.ReseedDatabase(h.dbPath); err != nil {
		common.InternalError(c, "Failed to reset database: "+err.Error())
		return
	}
	if err := seed.Seed(database.GetDB()); err != nil {
		common.InternalError(c, "Failed to seed database: "+err.Error())
		return
	}
	common.SuccessMessage(c, "Database reset successful")
}

// AnnouncementHandler 公告管理
type AnnouncementHandler struct {
	announcementService *service.AnnouncementService
}

func NewAnnouncementHandler(svc *service.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{announcementService: svc}
}

func (h *AnnouncementHandler) ListPublished(c *gin.Context) {
	list, err := h.announcementService.List()
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.Success(c, list)
}

func (h *AnnouncementHandler) Create(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	var input struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		IsPinned bool   `json:"is_pinned"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Title == "" || input.Content == "" {
		common.BadRequest(c, "Invalid request")
		return
	}
	a, err := h.announcementService.Create(user, input.Title, input.Content, input.IsPinned)
	if err != nil {
		common.Forbidden(c, err.Error())
		return
	}
	common.Success(c, a)
}

func (h *AnnouncementHandler) Update(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var input struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		IsPinned bool   `json:"is_pinned"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}
	if err := h.announcementService.Update(user, uint(id), input.Title, input.Content, input.IsPinned); err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.NotFound(c, err.Error())
		}
		return
	}
	common.SuccessMessage(c, "Announcement updated")
}

func (h *AnnouncementHandler) Delete(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.announcementService.Delete(user, uint(id)); err != nil {
		if err == service.ErrUnauthorized {
			common.Forbidden(c, err.Error())
		} else {
			common.NotFound(c, err.Error())
		}
		return
	}
	common.SuccessMessage(c, "Announcement deleted")
}

// ConfigHandler 系统配置管理
type ConfigHandler struct {
	configRepo *repository.SystemConfigRepository
}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{configRepo: repository.GetSystemConfigRepo()}
}

func (h *ConfigHandler) GetConfig(c *gin.Context) {
	configs, err := h.configRepo.List()
	if err != nil {
		common.InternalError(c, err.Error())
		return
	}

	// 返回默认值作为兜底
	result := map[string]string{
		"site_name":               "CloudNest",
		"allow_registration":      "true",
		"default_vps_expire_days": "30",
		"max_vps_per_user":        "10",
	}
	for _, cfg := range configs {
		result[cfg.Key] = cfg.Value
	}
	common.Success(c, result)
}

func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	var input struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Key == "" {
		common.BadRequest(c, "Invalid request")
		return
	}
	if err := h.configRepo.Upsert(input.Key, input.Value); err != nil {
		common.InternalError(c, err.Error())
		return
	}
	common.SuccessMessage(c, "Config updated")
}
