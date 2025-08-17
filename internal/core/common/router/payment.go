package router
import (
	"github.com/Khangvn20/FlyJourney_Backend/internal/controller"
	"github.com/gin-gonic/gin"
)

func PaymentRoutes(rg *gin.RouterGroup, paymentController *controller.PaymentController) {
	paymentRoutes := rg.Group("/payment")
	{
		paymentRoutes.POST("/momo", paymentController.CreatePayment)
		paymentRoutes.GET("/momo/success", paymentController.HandleMomoSuccess)
		paymentRoutes.POST("/momo/callback", paymentController.HandleMomoCallback)
	}
}