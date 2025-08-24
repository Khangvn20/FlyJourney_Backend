package dto 
import "time"
type CancellationEmailData struct {
    PNRCode            string
    BookingID          int64
    UserFullName       string
    ContactEmail       string
    ContactPhone       string
    CancellationReason string
    RefundAmount       float64
    OutboundFlight     *BookingEmailFlight
    InboundFlight      *BookingEmailFlight
    CancellationDate   time.Time
}
type DelayEmailData struct {
     PNRCode          string
    BookingID        int64
    UserFullName     string
    ContactEmail     string
    ContactPhone     string
    DelayReason      string
    OriginalTime     time.Time
    NewDepartureTime time.Time
    DelayDuration    string 
    NotificationTime time.Time
    OutboundFlight   *BookingEmailFlight
    InboundFlight    *BookingEmailFlight
}