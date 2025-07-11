package handlers

import (
	"net/http"
	"tuity/internal/adapters/driving/http/dto"
	"tuity/internal/core/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("validation_error", "Invalid request data", err.Error()))
		return
	}

	user, err := h.userService.CreateUser(req.Username, req.DisplayName)
	if err != nil {
		c.Error(err)
		return
	}

	response := dto.ToUserResponse(user)
	c.JSON(http.StatusCreated, response)
}

// GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userService.GetUser(userID)
	if err != nil {
		c.Error(err)
		return
	}

	response := dto.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// GET /users/:username
func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		c.Error(err)
		return
	}

	response := dto.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}
