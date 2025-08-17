package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
    "github.com/gin-gonic/gin"
	"encoding/json"
	"net/http"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
	"github.com/Khangvn20/FlyJourney_Backend/internal/infra/config"
)
type paymentService struct {
	momoConfig *config.MomoConfig
}
func NewPaymentService(momoConfig *config.MomoConfig) service.PaymentService {
    return &paymentService{
        momoConfig: momoConfig,
    }
}
func (s *paymentService) GenerateMomoSignature(req *request.MomoRequest) response.Response {
    // Đảm bảo sử dụng values từ request, không phải từ config
    rawData := fmt.Sprintf("accessKey=%s&amount=%s&extraData=%s&ipnUrl=%s&orderId=%s&orderInfo=%s&partnerCode=%s&redirectUrl=%s&requestId=%s&requestType=%s",
        req.AccessKey,    // Từ request
        req.Amount, 
        req.ExtraData, 
        req.IpnUrl,       // Từ request
        req.OrderId, 
        req.OrderInfo, 
        req.PartnerCode,  // Từ request
        req.RedirectUrl,  // Từ request (QUAN TRỌNG)
        req.RequestId, 
        req.RequestType)

    log.Printf("Raw data for signature: %s", rawData) // Debug log

    h := hmac.New(sha256.New, []byte(s.momoConfig.SecretKey))
    _, err := h.Write([]byte(rawData))
    if err != nil {
        return response.Response{
            Status:       false,
            ErrorCode:    error_code.SignatureError,
            ErrorMessage: fmt.Sprintf("Error generating signature: %v", err),
        }
    }

    signature := hex.EncodeToString(h.Sum(nil))
    log.Printf("Generated signature: %s", signature) // Debug log
    
    return response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Success",
        Data:         signature,
    }
}
func (s *paymentService) CreateMomoPayment(req *request.MomoRequest) response.Response {
       if req.PartnerCode == "" {
        req.PartnerCode = s.momoConfig.PartnerCode
    }
    if req.AccessKey == "" {
        req.AccessKey = s.momoConfig.AccessKey
    }
    if req.IpnUrl == "" {
        req.IpnUrl = s.momoConfig.IpnUrl
    }
    if req.RedirectUrl == "" {
        req.RedirectUrl = s.momoConfig.RedirectUrl
    }

    // Log request để debug
    log.Printf("Request before signature: %+v", req)
    signatureResponse := s.GenerateMomoSignature(req)
    if !signatureResponse.Status {
        return signatureResponse 
    }
    signature, ok := signatureResponse.Data.(string)
    if !ok {
        return response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Invalid signature format",
        }
    }
	req.Signature = signature
    payload, err := json.Marshal(req)
    if err != nil {
        return response.Response{
            Status:       false,
            ErrorCode:    "PAYLOAD_ERROR",
            ErrorMessage: fmt.Sprintf("Error marshalling payload: %v", err),
        }
    }
	resp, err := http.Post(s.momoConfig.Endpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return response.Response{
			Status:       false,
			ErrorCode:    error_code.InternalError,
			ErrorMessage: fmt.Sprintf("Error sending request: %v", err),
		}
	}
	defer resp.Body.Close()
    var momoResponse map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&momoResponse); err != nil {
        return response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: fmt.Sprintf("Error decoding response: %v", err),
        }
    }
   return response.Response{
       Status:       true,
       ErrorCode:    error_code.Success,
       ErrorMessage: error_code.SuccessErrMsg,
       Data:        momoResponse,
   }
}
// Private method - không cần khai báo trong interface
func (s *paymentService) handleSuccessfulPayment(req *request.MomoCallbackRequest) response.Response {
    // Log successful payment
    log.Printf("Payment successful - OrderID: %s, TransID: %s, Amount: %s", 
               req.OrderId, req.TransId, req.Amount)
    
    // TODO: Update booking status via BookingService
    // if s.bookingService != nil {
    //     s.bookingService.UpdateBookingStatus(req.OrderId, "paid")
    // }
    
    return response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Payment processed successfully",
        Data: map[string]interface{}{
            "orderId":      req.OrderId,
            "transId":      req.TransId,
            "amount":       req.Amount,
            "resultCode":   req.ResultCode,
            "message":      req.Message,
            "payType":      req.PayType,
            "responseTime": req.ResponseTime,
        },
    }
}

// Private method - không cần khai báo trong interface
func (s *paymentService) handleFailedPayment(req *request.MomoCallbackRequest, reason string) response.Response {
    // Log failed payment
    log.Printf("Payment failed - OrderID: %s, Reason: %s, ResultCode: %d", 
               req.OrderId, reason, req.ResultCode)
    
    // TODO: Update booking status via BookingService
    // if s.bookingService != nil {
    //     s.bookingService.UpdateBookingStatus(req.OrderId, "payment_failed")
    // }
    
    return response.Response{
        Status:       false,
        ErrorCode:    error_code.PaymentFailed,
        ErrorMessage: fmt.Sprintf("Payment failed: %s", reason),
        Data: map[string]interface{}{
            "orderId":      req.OrderId,
            "resultCode":   req.ResultCode,
            "message":      req.Message,
            "reason":       reason,
            "responseTime": req.ResponseTime,
        },
    }
}

// Private method - không cần khai báo trong interface
func (s *paymentService) verifyMomoCallbackSignature(req *request.MomoCallbackRequest) bool {
    rawData := fmt.Sprintf("accessKey=%s&amount=%s&extraData=%s&message=%s&orderId=%s&orderInfo=%s&orderType=%s&partnerCode=%s&payType=%s&requestId=%s&responseTime=%s&resultCode=%d&transId=%s",
        s.momoConfig.AccessKey, req.Amount, req.ExtraData, req.Message,
        req.OrderId, req.OrderInfo, req.OrderType, req.PartnerCode,
        req.PayType, req.RequestId, req.ResponseTime, req.ResultCode, req.TransId)

    h := hmac.New(sha256.New, []byte(s.momoConfig.SecretKey))
    h.Write([]byte(rawData))
    expectedSignature := hex.EncodeToString(h.Sum(nil))

    return expectedSignature == req.Signature
}
func (s *paymentService) HandleMomoCallback(req *request.MomoCallbackRequest) response.Response {
    // Verify callback signature
    if !s.verifyMomoCallbackSignature(req) {
        return response.Response{
            Status:       false,
            ErrorCode:    error_code.SignatureError,
            ErrorMessage: "Invalid callback signature",
        }
    }

    // Process payment result based on resultCode
    switch req.ResultCode {
    case 0:
        // Payment successful
        return s.handleSuccessfulPayment(req)
    case 1006:
        // Transaction timeout
        return s.handleFailedPayment(req, "Transaction timeout")
    case 1007:
        // Insufficient balance
        return s.handleFailedPayment(req, "Insufficient balance")
    case 1009:
        // Transaction cancelled by user
        return s.handleFailedPayment(req, "Transaction cancelled by user")
    case 1010:
        // User rejected transaction
        return s.handleFailedPayment(req, "User rejected transaction")
    default:
        // Other error
        return s.handleFailedPayment(req, fmt.Sprintf("Payment failed with code: %d", req.ResultCode))
    }
}
// Trong file internal/core/service/payment.go
func (s *paymentService) HandleMomoSuccess(ctx *gin.Context) response.Response {
    // Extract query parameters
    partnerCode := ctx.Query("partnerCode")
    orderId := ctx.Query("orderId")
    requestId := ctx.Query("requestId")
    amount := ctx.Query("amount")
    orderInfo := ctx.Query("orderInfo")
    orderType := ctx.Query("orderType")
    transId := ctx.Query("transId")
    resultCode := ctx.Query("resultCode")
    message := ctx.Query("message")
    payType := ctx.Query("payType")
    responseTime := ctx.Query("responseTime")
    extraData := ctx.Query("extraData")
    signature := ctx.Query("signature")

    // Validate required fields
    if partnerCode == "" || orderId == "" || resultCode == "" {
        return response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Missing required parameters",
        }
    }

    // Log success redirect
    log.Printf("MoMo redirect success - OrderID: %s, ResultCode: %s, Message: %s", 
               orderId, resultCode, message)

    if resultCode == "0" {
  
        return response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            ErrorMessage: "Payment completed successfully",
            Data: map[string]interface{}{
                "orderId":      orderId,
                "requestId":    requestId,
                "transId":      transId,
                "amount":       amount,
                "orderInfo":    orderInfo,
                "orderType":    orderType,  // Sử dụng orderType ở đây
                "resultCode":   resultCode,
                "message":      message,
                "payType":      payType,
                "responseTime": responseTime,
                "extraData":    extraData,
                "signature":    signature,
            },
        }
    } else {
     
        return response.Response{
            Status:       false,
            ErrorCode:    error_code.PaymentFailed,
            ErrorMessage: fmt.Sprintf("Payment failed: %s", message),
            Data: map[string]interface{}{
                "orderId":    orderId,
                "resultCode": resultCode,
                "message":    message,
            },
        }
    }
}