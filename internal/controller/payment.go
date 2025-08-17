package controller 
import (
    "net/http"
	"log"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
    "github.com/gin-gonic/gin"
)
type PaymentController struct {
	paymentService service.PaymentService
}
func NewPaymentController(paymentService service.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}
func (c *PaymentController) CreatePayment(ctx *gin.Context) {
    var req request.MomoRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        log.Printf("Invalid payment request: %v", err)
        ctx.JSON(http.StatusBadRequest, response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid payment request format",
        })
        return
    }
    if req.Amount == "" || req.OrderId == "" {
        ctx.JSON(http.StatusBadRequest, response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Missing required fields",
        })
        return
    }
    
    log.Printf("Creating MoMo payment for order: %s, amount: %s", req.OrderId, req.Amount)
    
    result := c.paymentService.CreateMomoPayment(&req)
    
    if result.Status {
        log.Printf("Payment created successfully for order: %s", req.OrderId)
        ctx.JSON(http.StatusOK, result)
    } else {
        log.Printf("Payment creation failed for order: %s, error: %s", req.OrderId, result.ErrorMessage)
        ctx.JSON(http.StatusBadRequest, result)
    }
}

func (c *PaymentController) HandleMomoCallback(ctx *gin.Context) {
    var req request.MomoCallbackRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        log.Printf("Invalid MoMo callback: %v", err)
        ctx.JSON(http.StatusBadRequest, response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid callback format",
        })
        return
    }
    
    log.Printf("Received MoMo callback for order: %s, result: %d", req.OrderId, req.ResultCode)
    
    result := c.paymentService.HandleMomoCallback(&req)
    
    // Always return 200 for MoMo callback acknowledgment
    ctx.JSON(http.StatusOK, result)
}

func (c *PaymentController) HandleMomoSuccess(ctx *gin.Context) {
    result := c.paymentService.HandleMomoSuccess(ctx)
    
    // Trả về HTML hoặc JSON tùy theo yêu cầu
    ctx.JSON(http.StatusOK, result)
    
    // Hoặc redirect về frontend nếu cần:
    // if result.Status {
    //     ctx.Redirect(http.StatusFound, "http://localhost:3000/payment/success")
    // } else {
    //     ctx.Redirect(http.StatusFound, "http://localhost:3000/payment/failed")
    // }
}