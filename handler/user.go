package handler

import (
	"crowdfunding/auth"
	"crowdfunding/helper"
	"crowdfunding/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// init variable to store request data
	var input user.RegisterUserInput
	// store request data to input variable
	err := c.ShouldBindJSON(&input)
	if err != nil {
		// format an errors
		errors := helper.FormatValidationError(err)
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

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Login user failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loggedInUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login user failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login user failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, token)
	response := helper.APIResponse("User login successfully", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmail(c *gin.Context) {
	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Check email failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Check email failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_available": isEmailAvailable}
	metaMessage := "Email has been registered"
	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false, "error": err.Error()}
		response := helper.APIResponse("Upload avatar failed", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// get user data from middleware
	usr := c.MustGet("currentUser").(user.User)
	userID := usr.ID
	userIDSplit := strings.Split(userID, "-")[4]
	currentImage := usr.AvatarFileName

	// delete current images
	if currentImage != "" {
		if err := os.Remove(currentImage); err != nil {
			fmt.Println("user doesn't have images.", err.Error())
		}
	}

	// rename images
	extension := filepath.Ext(file.Filename)
	randomText := strings.Split(uuid.New().String(), "-")[4]
	randomText2 := strings.Split(uuid.New().String(), "-")[4]
	//fileName := strings.TrimSuffix(file.Filename, extension)
	//filePath := "images/" + fileName + "-" + randomText + "-" + userIDSplit + extension
	filePath := "images/user-avatar/" + randomText + randomText2 + userIDSplit + extension
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		data := gin.H{"is_uploaded": false, "error": err.Error()}
		response := helper.APIResponse("Upload avatar failed", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// save image location to db
	_, err = h.userService.SaveAvatar(userID, filePath)
	if err != nil {
		data := gin.H{"is_uploaded": false, "error": err.Error()}
		response := helper.APIResponse("Upload avatar failed", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Upload avatar successfully", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
