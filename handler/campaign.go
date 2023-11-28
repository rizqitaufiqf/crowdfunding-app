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
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := campaign.FormatCampaigns(campaigns)
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}
