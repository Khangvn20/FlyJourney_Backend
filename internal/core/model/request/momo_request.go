package request

type MomoRequest struct {
    BookingID   string  `json:"booking_id" binding:"required"`
	PartnerCode string `json:"partnerCode"`
	AccessKey   string `json:"accessKey"`
	RequestId   string `json:"requestId"`
	Amount      string `json:"amount"`
	OrderId     string `json:"orderId"`
	OrderInfo   string `json:"orderInfo"`
	RedirectUrl string `json:"redirectUrl"`
	IpnUrl      string `json:"ipnUrl"`
	ExtraData   string `json:"extraData"`
	RequestType string `json:"requestType"`
	Signature   string `json:"signature"`
}
type MomoCallbackRequest struct {
    PartnerCode     string `json:"partnerCode" binding:"required"`
    OrderId         string `json:"orderId" binding:"required"`
    RequestId       string `json:"requestId" binding:"required"`
    Amount          int64 `json:"amount" binding:"required"`
    OrderInfo       string `json:"orderInfo"`
    OrderType       string `json:"orderType"`
    TransId         int64 `json:"transId"`
    ResultCode      int    `json:"resultCode"`
    Message         string `json:"message"`
    PayType         string `json:"payType"`
    ResponseTime    int64 `json:"responseTime"`
    ExtraData       string `json:"extraData"`
    Signature       string `json:"signature" binding:"required"`
}