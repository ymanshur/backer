package handler

import (
	"backer/campaign"
	"backer/helper"
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
