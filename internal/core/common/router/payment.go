package router
import (
	"github.com/Khangvn20/FlyJourney_Backend/internal/controller"
	"github.com/gin-gonic/gin"
)

func PaymentRoutes(rg *gin.RouterGroup, paymentController *controller.PaymentController,authMiddleware gin.HandlerFunc) {
	paymentRoutes := rg.Group("/payment")
	{
		paymentRoutes.POST("/momo",authMiddleware ,paymentController.CreatePayment)

	}
}