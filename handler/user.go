package handler

import (
	"backer/auth"
	"backer/helper"
	"backer/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(ctx *gin.Context) {
	/**
	 * 1. Get client input
	 * 2. Map the input into RegisterUserInput struct (DTO)
	 * 3. Pass the DTO into service
	 */

	var input user.RegisterUserInput

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Register account failed", http.StatusInternalServerError, "error", errorMessage)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Register account failed", http.StatusInternalServerError, "error", errorMessage)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(ctx *gin.Context) {
	/**
	 * 1. Client input 'email' and 'password', then handler catch it request data
	 * 2. Handler mapping the input into struct (DTO)
	 * 3. Passing the mapped struct into service
	 * 4. Service find email payload with data in users table (database)
	 * 5. Matching payload password with expect password
	 */

	var input user.LoginInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login failed", http.StatusInternalServerError, "error", errorMessage)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, token)

	response := helper.APIResponse("Successfully logged in", http.StatusOK, "success", formatter)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(ctx *gin.Context) {
	/**
	 * 1. Get input 'email' from client
	 * 2. Mapping input into struct
	 * 3. Pass email value to service
	 * 4. Service using repository to check email availability
	 */

	var input user.CheckEmailInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Email checking failed", http.StatusInternalServerError, "error", errorMessage)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	data := gin.H{"is_available": isEmailAvailable}

	var metaMessage string

	if isEmailAvailable {
		metaMessage = "Email is available"
	} else {
		metaMessage = "Email has been registered"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(ctx *gin.Context) {
	/**
	 * 1. Get 'avatar' input form from client
	 * 2. Store the image into "images/"
	 * 3. Save input into database via service
	 *      i. Get id of client from JWT token
	 *      ii. Get user data with the id, then
	 *      iii. Save user with uploaded avatar (only the path)
	 */

	file, err := ctx.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// Should be got from JWT token
	userID := 1

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	if ctx.SaveUploadedFile(file, path); err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload avatar image", http.StatusInternalServerError, "error", data)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if _, err := h.userService.SaveAvatar(userID, path); err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload avatar image", http.StatusInternalServerError, "error", data)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	data := gin.H{"is_uploaded": true}

	response := helper.APIResponse("Avatar successfully uploaded", http.StatusOK, "success", data)
	ctx.JSON(http.StatusOK, response)
}
