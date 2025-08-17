package error_code

type ErrorCode string

// error code
const (
	Success        ErrorCode = "SUCCESS"
	InvalidRequest ErrorCode = "INVALID_REQUEST"
	DuplicateUser  ErrorCode = "DUPLICATE_USER"
	InternalError  ErrorCode = "INTERNAL_ERROR"
	Forbidden      ErrorCode = "FORBIDDEN"
	InvalidID      ErrorCode = "INVALID_ID"
	ResourceLocked ErrorCode = "RESOURCE_LOCKED"
	NOTFOUND       ErrorCode = "NOT_FOUND"
	InsufficientSeats ErrorCode = "INSUFFICIENT_SEATS"
	NoAvailableSeats ErrorCode = "NO_AVAILABLE_SEATS"
	RateLimitExceeded ErrorCode = "RATE_LIMIT_EXCEEDED"
	SignatureError     ErrorCode = "SIGNATURE_ERROR"
	PaymentFailed      ErrorCode = "PAYMENT_FAILED"
)

// error message
const (
	SuccessErrMsg        = "success"
	InternalErrMsg       = "internal error"
	InvalidRequestErrMsg = "invalid request"
)
