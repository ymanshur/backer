package handler

import (
	"backer/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(ctx *gin.Context) {
	/**
	 * RegisterUser handler:
	 * 1. Get client input
	 * 2. Map the input into RegisterUserInput struct (DTO)
	 * 3. Pass the DTO into service
	 */

	var input user.RegisterUserInput

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
	}

	ctx.JSON(http.StatusOK, newUser)
}
