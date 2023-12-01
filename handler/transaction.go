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
		response := helper.APIResponse("Failed to get campaign transaction", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	usr := c.MustGet("currentUser").(user.User)
	input.User = usr

	trans, err := h.service.GetTransactionByCampaignID(input)
	if err != nil {
		if strings.Contains(err.Error(), "invalid user") {
			response := helper.APIResponse("Failed to get campaign transaction", http.StatusUnauthorized, "error", err.Error())
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		response := helper.APIResponse("Failed to get campaign transaction", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign transactions", http.StatusOK, "success", transactions.FormatCampaignTransactions(trans))
	c.JSON(http.StatusOK, response)
}
