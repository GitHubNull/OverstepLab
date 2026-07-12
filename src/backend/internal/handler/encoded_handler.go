package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oversteplab/oversteplab/internal/common"
	cryptopkg "github.com/oversteplab/oversteplab/internal/crypto"
	"github.com/oversteplab/oversteplab/internal/middleware"
	"github.com/oversteplab/oversteplab/internal/service"
)

// EncodedHandler wraps existing handlers with encoding/encryption decoding support.
// It decodes the path parameter ":encodedId" using the encoding type specified in
// the X-Encoding-Type request header, then delegates to the actual service layer.
type EncodedHandler struct {
	vpsSvc     *service.VPSService
	userSvc    *service.UserService
	orderSvc   *service.OrderService
	ticketSvc  *service.TicketService
	apiKeySvc  *service.APIKeyService
	companySvc *service.CompanyService
}

func NewEncodedHandler(
	vpsSvc *service.VPSService,
	userSvc *service.UserService,
	orderSvc *service.OrderService,
	ticketSvc *service.TicketService,
	apiKeySvc *service.APIKeyService,
	companySvc *service.CompanyService,
) *EncodedHandler {
	return &EncodedHandler{
		vpsSvc:     vpsSvc,
		userSvc:    userSvc,
		orderSvc:   orderSvc,
		ticketSvc:  ticketSvc,
		apiKeySvc:  apiKeySvc,
		companySvc: companySvc,
	}
}

// getEncodingType extracts the encoding type from the X-Encoding-Type header.
func getEncodingType(c *gin.Context) string {
	encType := c.GetHeader("X-Encoding-Type")
	if encType == "" {
		encType = "base64" // default
	}
	return encType
}

// decodeEncodedID decodes the ":encodedId" path parameter or "?v=" query parameter
// using the specified encoding type. Query parameter takes priority for complex
// encodings (like SM4/AES) where base64 `+` chars may be URL-decoded in path params.
func decodeEncodedID(c *gin.Context, encType string) (uint, error) {
	encodedID := c.Query("v")
	if encodedID == "" {
		encodedID = c.Param("encodedId")
	}
	if encodedID == "" {
		return 0, service.ErrVPSNotFound
	}

	decoded, err := cryptopkg.DecodeParam(encodedID, encType)
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseUint(decoded, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

// handleDecodeError returns appropriate HTTP error for decoding failures.
func handleDecodeError(c *gin.Context, err error) {
	common.BadRequest(c, "Decode error: "+err.Error()+" - check X-Encoding-Type header and encoded parameter format")
}

// ---- VPS Encoded Endpoints (H-01, V-01, V-03, C-01 encoding variants) ----

// GetVPSByEncodedID handles GET /api/v1/encoded/vps/:encodedId
func (h *EncodedHandler) GetVPSByEncodedID(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	encType := getEncodingType(c)

	id, err := decodeEncodedID(c, encType)
	if err != nil {
		handleDecodeError(c, err)
		return
	}

	vps, err := h.vpsSvc.GetDetail(user, id)
	if err != nil {
		common.NotFound(c, "VPS not found")
		return
	}
	common.Success(c, vps)
}

// StartVPSEncoded handles POST /api/v1/encoded/vps/:encodedId/start
func (h *EncodedHandler) StartVPSEncoded(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	encType := getEncodingType(c)

	id, err := decodeEncodedID(c, encType)
	if err != nil {
		handleDecodeError(c, err)
		return
	}

	if err := h.vpsSvc.StartVPS(user, id); err != nil {
		common.Forbidden(c, err.Error())
		return
	}
	common.SuccessMessage(c, "VPS started")
}

// StopVPSEncoded handles POST /api/v1/encoded/vps/:encodedId/stop
func (h *EncodedHandler) StopVPSEncoded(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	encType := getEncodingType(c)

	id, err := decodeEncodedID(c, encType)
	if err != nil {
		handleDecodeError(c, err)
		return
	}

	if err := h.vpsSvc.StopVPS(user, id); err != nil {
		common.Forbidden(c, err.Error())
		return
	}
	common.SuccessMessage(c, "VPS stopped")
}

// ReinstallVPSEncoded handles POST /api/v1/encoded/vps/:encodedId/reinstall
func (h *EncodedHandler) ReinstallVPSEncoded(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	encType := getEncodingType(c)

	id, err := decodeEncodedID(c, encType)
	if err != nil {
		handleDecodeError(c, err)
		return
	}

	var input struct {
		OSImage string `json:"os_image"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}

	if err := h.vpsSvc.ReinstallVPS(user, id, input.OSImage); err != nil {
		common.Forbidden(c, err.Error())
		return
	}
	common.SuccessMessage(c, "VPS reinstalled")
}

// ---- User Encoded Endpoint (H-02 encoding variant) ----

// GetUserByEncodedID handles GET /api/v1/encoded/users/:encodedId
func (h *EncodedHandler) GetUserByEncodedID(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	encType := getEncodingType(c)

	id, err := decodeEncodedID(c, encType)
	if err != nil {
		handleDecodeError(c, err)
		return
	}

	user, err := h.userSvc.GetUserByID(currentUser, id)
	if err != nil {
		common.NotFound(c, "User not found")
		return
	}
	common.Success(c, user)
}

// ---- Order Encoded Endpoint (H-03 encoding variant) ----

// GetOrderByEncodedID handles GET /api/v1/encoded/orders/:encodedId
func (h *EncodedHandler) GetOrderByEncodedID(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	encType := getEncodingType(c)

	id, err := decodeEncodedID(c, encType)
	if err != nil {
		handleDecodeError(c, err)
		return
	}

	order, err := h.orderSvc.GetDetail(user, id)
	if err != nil {
		common.NotFound(c, "Order not found")
		return
	}
	common.Success(c, order)
}

// ---- Ticket Encoded Endpoints (H-04 encoding variant) ----

// GetTicketByEncodedID handles GET /api/v1/encoded/tickets/:encodedId
func (h *EncodedHandler) GetTicketByEncodedID(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	encType := getEncodingType(c)

	id, err := decodeEncodedID(c, encType)
	if err != nil {
		handleDecodeError(c, err)
		return
	}

	ticket, err := h.ticketSvc.GetDetail(user, id)
	if err != nil {
		common.NotFound(c, "Ticket not found")
		return
	}
	common.Success(c, ticket)
}

// ---- API Key Encoded Endpoint (H-05 encoding variant) ----

// DeleteAPIKeyEncoded handles DELETE /api/v1/encoded/apikeys/:encodedId
func (h *EncodedHandler) DeleteAPIKeyEncoded(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	encType := getEncodingType(c)

	id, err := decodeEncodedID(c, encType)
	if err != nil {
		handleDecodeError(c, err)
		return
	}

	if err := h.apiKeySvc.Delete(user, id); err != nil {
		common.Forbidden(c, err.Error())
		return
	}
	common.SuccessMessage(c, "API Key deleted")
}

// ---- Company Encoded Endpoints (V-02, V-05, C-02 encoding variants) ----

// AddMemberEncoded handles POST /api/v1/encoded/company/members (V-02/C-02: body params encoded)
func (h *EncodedHandler) AddMemberEncoded(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}

	// Body fields are already decoded by EncodingMiddleware

	member, err := h.companySvc.AddMember(user, &service.AddMemberInput{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
		Role:     input.Role,
	})
	if err != nil {
		common.Forbidden(c, err.Error())
		return
	}
	common.Success(c, member)
}

// ChangeRoleEncoded handles PUT /api/v1/encoded/company/members/:encodedId/role (V-05 encoding variant)
func (h *EncodedHandler) ChangeRoleEncoded(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	encType := getEncodingType(c)

	id, err := decodeEncodedID(c, encType)
	if err != nil {
		handleDecodeError(c, err)
		return
	}

	var input struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}

	if err := h.companySvc.ChangeRole(user, id, input.Role); err != nil {
		common.Forbidden(c, err.Error())
		return
	}
	common.SuccessMessage(c, "Role updated")
}

// ---- Crypto Utility Endpoints (for challenge learning) ----

// EncodeValue handles POST /api/v1/crypto/encode
func (h *EncodedHandler) EncodeValue(c *gin.Context) {
	var input struct {
		Value        string `json:"value"`
		EncodingType string `json:"encoding_type"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Value == "" {
		common.BadRequest(c, "Invalid request: value and encoding_type required")
		return
	}

	var result string
	var err error

	switch input.EncodingType {
	case "base64":
		result = cryptopkg.Base64Encode([]byte(input.Value))
	case "base32":
		result = cryptopkg.Base32Encode([]byte(input.Value))
	case "base58":
		result = cryptopkg.Base58Encode([]byte(input.Value))
	case "base85":
		result = cryptopkg.Base85Encode([]byte(input.Value))
	case "custom_base64":
		result = cryptopkg.CustomBase64Encode([]byte(input.Value))
	case "custom_base32":
		result = cryptopkg.CustomBase32Encode([]byte(input.Value))
	case "caesar":
		result = cryptopkg.CaesarEncodeWithShift(input.Value, 3)
	case "aes":
		result, err = cryptopkg.AESEncrypt([]byte(input.Value))
	case "rsa":
		result, err = cryptopkg.RSAEncrypt([]byte(input.Value))
	case "sm2":
		result, err = cryptopkg.SM2Encrypt([]byte(input.Value))
	case "sm4":
		result, err = cryptopkg.SM4Encrypt([]byte(input.Value))
	case "signed":
		result = cryptopkg.EncodeSignedParam(input.Value)
	case "hash":
		// E-09: hash is a header signature, not a per-value encoding
		// For the encode utility, just return a demo hash of the value
		result = cryptopkg.ComputeHashSign(input.Value)
	default:
		common.BadRequest(c, "Unsupported encoding type: "+input.EncodingType)
		return
	}

	if err != nil {
		common.InternalError(c, "Encode error: "+err.Error())
		return
	}

	common.Success(c, gin.H{
		"encoding_type": input.EncodingType,
		"original":      input.Value,
		"encoded":       result,
	})
}

// DecodeValue handles POST /api/v1/crypto/decode
func (h *EncodedHandler) DecodeValue(c *gin.Context) {
	var input struct {
		Value        string `json:"value"`
		EncodingType string `json:"encoding_type"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Value == "" {
		common.BadRequest(c, "Invalid request: value and encoding_type required")
		return
	}

	decoded, err := cryptopkg.DecodeByType(input.Value, input.EncodingType)
	if err != nil {
		common.BadRequest(c, "Decode error: "+err.Error())
		return
	}

	common.Success(c, gin.H{
		"encoding_type": input.EncodingType,
		"encoded":       input.Value,
		"decoded":       string(decoded),
	})
}

// GetCryptoKeys returns public crypto keys for challenge participants.
func (h *EncodedHandler) GetCryptoKeys(c *gin.Context) {
	common.Success(c, gin.H{
		"aes_key_base64":  cryptopkg.GetAESKeyBase64(),
		"hmac_key_base64": cryptopkg.GetHMACKeyBase64(),
		"sm4_key_base64":  cryptopkg.GetSM4KeyBase64(),
		"rsa_public_key":  cryptopkg.GetRSAPublicKeyPEM(),
		"sm2_public_key":  cryptopkg.GetSM2PublicKeyHex(),
		"hash_salt":       string(cryptopkg.GetHashSalt()),
	})
}
