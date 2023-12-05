package handler

import (
	"crowdfunding/helper"
	"crowdfunding/transactions"
	"crowdfunding/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type transactionHandler struct {
	service transactions.Service
}

func NewTansactionHandler(service transactions.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransaction(c *gin.Context) {
	var input transactions.GetCampaignTransactionDTO
	if err := c.ShouldBindUri(&input); err != nil {
		response := helper.APIResponse("Failed to get campaign transactions", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	usr := c.MustGet("currentUser").(user.User)
	input.User = usr

	trans, err := h.service.GetTransactionsByCampaignID(input)
	if err != nil {
		if strings.Contains(err.Error(), "invalid user") {
			response := helper.APIResponse("Failed to get campaign transactions", http.StatusUnauthorized, "error", err.Error())
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		response := helper.APIResponse("Failed to get campaign transactions", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign transactions", http.StatusOK, "success", transactions.FormatCampaignTransactions(trans))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	usr := c.MustGet("currentUser").(user.User)
	userID := usr.ID

	trans, err := h.service.GetTransactionsByUserID(userID)
	if err != nil {
		response := helper.APIResponse("Failed to get user transactions", http.StatusUnauthorized, "error", err.Error())
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	response := helper.APIResponse("Campaign transactions", http.StatusOK, "success", transactions.FormatUserTransactions(trans))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transactions.CreateTransactionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create transaction", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	usr := c.MustGet("currentUser").(user.User)
	input.User = usr

	trans, err := h.service.CreateTransaction(input)
	if err != nil {
		response := helper.APIResponse("Failed to create transaction", http.StatusUnauthorized, "error", err.Error())
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	response := helper.APIResponse("Create transaction successfully", http.StatusCreated, "success", transactions.FormatCreateTransaction(trans))
	c.JSON(http.StatusCreated, response)
}
