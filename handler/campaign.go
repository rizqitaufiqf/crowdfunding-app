package handler

import (
	"crowdfunding/campaign"
	"crowdfunding/helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID := c.Query("user_id")
	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil && !strings.Contains(err.Error(), "invalid input syntax for type uuid") {
		response := helper.APIResponse("Failed to get campaigns", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := campaign.FormatCampaigns(campaigns)
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaignDetail(c *gin.Context) {
	var input campaign.GetCampaignDTO

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail campaign", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail campaign", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}
