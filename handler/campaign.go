package handler

import (
	"backer/campaign"
	"backer/helper"
	"backer/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse(
			"Error to get campaigns",
			http.StatusInternalServerError,
			"error",
			errorMessage,
		)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse(
		"List of campaigns",
		http.StatusOK,
		"success",
		campaign.FormatCampaigns(campaigns),
	)
	ctx.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(ctx *gin.Context) {
	/**
	 * 1. Handler: mapping 'id' from url into struct input, pass the input into service, format the result for response
	 * 2. Service: using 'id' at the input for repo param
	 * 3. Repository: get campaign by id
	 */

	var input campaign.GetCampaignInput

	if err := ctx.ShouldBindUri(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"Failed to get detail campaign",
			http.StatusBadRequest,
			"error",
			errorMessage,
		)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse(
			"Failed to get detail campaign",
			http.StatusInternalServerError,
			"error",
			errorMessage,
		)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse(
		"Campaign detail",
		http.StatusOK,
		"success",
		campaign.FormatCampaignDetail(campaignDetail),
	)
	ctx.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(ctx *gin.Context) {
	var input campaign.CreateCampaignInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"Failed to create campaign",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(user.User)

	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse(
			"Failed to create campaign",
			http.StatusInternalServerError,
			"error",
			errorMessage,
		)
		ctx.JSON(http.StatusInternalServerError, response)
	}

	response := helper.APIResponse(
		"Campaign successfully created",
		http.StatusOK,
		"success",
		campaign.FormatCampaign(newCampaign),
	)
	ctx.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign(ctx *gin.Context) {
	var inputID campaign.GetCampaignInput

	if err := ctx.ShouldBindUri(&inputID); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"Failed to update campaign",
			http.StatusBadRequest,
			"error",
			errorMessage,
		)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	if err := ctx.ShouldBindJSON(&inputData); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"Failed to update campaign",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
		)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(user.User)

	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse(
			"Failed to update campaign",
			http.StatusBadRequest,
			"error",
			errorMessage,
		)
		ctx.JSON(http.StatusBadRequest, response)
	}

	response := helper.APIResponse(
		"Campaign successfully updated",
		http.StatusOK,
		"success",
		campaign.FormatCampaign(updatedCampaign),
	)
	ctx.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadCampaignImage(ctx *gin.Context) {
	var input campaign.CreateCampaignImageInput

	if err := ctx.ShouldBind(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse(
			"Failed to upload campaign image",
			http.StatusBadRequest,
			"error",
			errorMessage,
		)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload campaign image", http.StatusInternalServerError, "error", data)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(user.User)

	path := fmt.Sprintf("campaign-images/%d-%s", currentUser.ID, file.Filename)

	if ctx.SaveUploadedFile(file, path); err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload campaign image", http.StatusInternalServerError, "error", data)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	input.User = currentUser

	if _, err := h.service.CreateCampaignImage(input, path); err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}

	response := helper.APIResponse("Campaign image successfully uploaded", http.StatusOK, "success", data)
	ctx.JSON(http.StatusOK, response)
}
