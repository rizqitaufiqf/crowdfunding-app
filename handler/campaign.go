package handler

import (
	"crowdfunding/campaign"
	"crowdfunding/helper"
	"crowdfunding/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
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
		if strings.Contains(err.Error(), "found") {
			response := helper.APIResponse("Failed to get detail campaign", http.StatusNotFound, "error", err.Error())
			c.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponse("Failed to get detail campaign", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignDTO

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Create campaign failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	usr := c.MustGet("currentUser").(user.User)
	input.User = usr

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Create campaign failed", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Create campaign successfully", http.StatusCreated, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusCreated, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDTO
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignDTO
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	usr := c.MustGet("currentUser").(user.User)
	inputData.User = usr

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		if strings.Contains(err.Error(), "invalid user") {
			response := helper.APIResponse("Failed to update campaign", http.StatusUnauthorized, "error", err.Error())
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Update campaign successfully", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageDTO
	if err := c.ShouldBind(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"is_uploaded": false, "errors": errors}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	file, err := c.FormFile("campaign_image")
	if err != nil {
		data := gin.H{"is_uploaded": false, "errors": []string{err.Error()}}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	usr := c.MustGet("currentUser").(user.User)
	userID := usr.ID
	userIDSplit := strings.Split(userID, "-")[4]

	extension := filepath.Ext(file.Filename)
	randomText := strings.Split(uuid.New().String(), "-")[4]
	randomText2 := strings.Split(uuid.New().String(), "-")[4]
	filePath := fmt.Sprintf("images/campaigns/%s/%s/%s%s%s%s", userID, input.CampaignID, randomText, randomText2, userIDSplit, extension)

	input.User = usr
	if _, err := h.service.SaveCampaignImage(input, filePath); err != nil {
		if strings.Contains(err.Error(), "invalid user") {
			data := gin.H{"is_uploaded": false, "errors": []string{err.Error()}}
			response := helper.APIResponse("Failed to upload campaign image", http.StatusUnauthorized, "error", data)
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		data := gin.H{"is_uploaded": false, "errors": []string{err.Error()}}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		data := gin.H{"is_uploaded": false, "errors": []string{err.Error()}}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image successfully uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
