package controller

import (
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
	"github.com/gin-gonic/gin"
	"log"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status":       false,
			"errorCode":    "INVALID_REQUEST",
			"errorMessage": err.Error(),
			"data":         nil,
		})
		return
	}

	result := c.userService.Register(&req)

	var statusCode int
	if result.Status {
		statusCode = 200
	} else {
		statusCode = 400
	}
	ctx.JSON(statusCode, result)
}

func (c *UserController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(400, gin.H{
			"status":       false,
			"errorCode":    "INVALID_REQUEST",
			"errorMessage": err.Error(),
			"data":         nil,
		})
		return
	}

	log.Printf("Processing login request for email: %s", req.Email)

	result := c.userService.Login(&req)

	var statusCode int
	if result.Status {
		statusCode = 200
		log.Printf("Login successful for email: %s", req.Email)
	} else {
		statusCode = 400
		log.Printf("Login failed for email: %s with error: %s", req.Email, result.ErrorMessage)
	}

	ctx.JSON(statusCode, result)
}
func (c *UserController) ConfirmRegister(ctx *gin.Context) {
    var req request.ConfirmRegisterRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, gin.H{"status": false, "errorMessage": err.Error()})
        return
    }
    result := c.userService.ConfirmRegister(&req)
    ctx.JSON(200, result)
}
func (c *UserController) ResetPassword(ctx *gin.Context) {
	var req request.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"status": false, "errorMessage": err.Error()})
		return
	}
	result := c.userService.ResetPassword(&req)
	ctx.JSON(200, result)
}
func (c *UserController) ConfirmResetPassword(ctx *gin.Context) {
	var req request.ConfirmResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"status": false, "errorMessage": err.Error()})
		return
	}	
	result := c.userService.ConfirmResetPassword(&req)
	ctx.JSON(200, result)
}
func (c *UserController) Logout(ctx *gin.Context) {
    tokenString := ctx.GetHeader("Authorization")
    if tokenString == "" {
        ctx.JSON(400, gin.H{
            "status":       false,
            "errorCode":    "INVALID_REQUEST",
            "errorMessage": "Authorization header is required",
            "data":         nil,
        })
        return
    }
    result := c.userService.Logout(tokenString)
    var statusCode int
    if result.Status {
        statusCode = 200
    } else {
        statusCode = 400
    }
    ctx.JSON(statusCode, result)
}
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(400, gin.H{
			"status":       false,
			"errorCode":    "INVALID_REQUEST",
			"errorMessage": "User ID not found in context",
			"data":         nil,
		})
		return
	}

	result := c.userService.GetUserInfo(userID.(int))

	var statusCode int
	if result.Status {
		statusCode = 200
	} else {
		statusCode = 400
	}

	ctx.JSON(statusCode, result)
}