package handler

import (
	"crowdfunding/helper"
	"crowdfunding/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// init variable to store request data
	var input user.RegisterUserInput

	// store request data to input variable
	err := c.ShouldBindJSON(&input)
	if err != nil {
		// store errors to slice
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Error())
		}
		// store slice of errors to object errors
		errorMessage := gin.H{"errors": errors}
		// construct response
		response := helper.APIResponse("Register user failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// register user
	usr, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Register user failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// format user response
	formatter := user.FormatUser(usr, "")
	// construct response
	response := helper.APIResponse("User has been registered", http.StatusCreated, "success", formatter)
	c.JSON(http.StatusOK, response)
}
