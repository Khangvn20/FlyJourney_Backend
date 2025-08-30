package request

type SendMessageRequest struct {
    ConversationID string `json:"conversation_id,omitempty"`
    SenderID       int64  `json:"sender_id" binding:"required"`
    ReceiverID     int64  `json:"receiver_id" binding:"required"`
    SenderType     string `json:"sender_type" binding:"required"` // "admin" hoặc "user"
    Content        string `json:"content" binding:"required"`
    AttachmentURL  string `json:"attachment_url,omitempty"`
    MessageType    string `json:"message_type" binding:"required"`
}

type GetMessagesRequest struct {
    ConversationID string `json:"conversation_id" binding:"required"`
    Limit          int    `json:"limit,omitempty"`
    Offset         int    `json:"offset,omitempty"`
}

type GetConversationsRequest struct {
    UserID  int64 `json:"user_id" binding:"required"`
    IsAdmin bool  `json:"is_admin"`
    Limit   int   `json:"limit,omitempty"`
    Offset  int   `json:"offset,omitempty"`
}

type MarkAsReadRequest struct {
    MessageID int64 `json:"message_id" binding:"required"`
    UserID    int64 `json:"user_id,omitempty"`
}

type GetUnreadCountRequest struct {
    UserID  int64 `json:"user_id" binding:"required"`
    IsAdmin bool  `json:"is_admin"`
}
type GetConversationHistoryRequest struct {
    ConversationID string `json:"conversation_id" binding:"required"`
    Limit          int    `json:"limit,omitempty"`
    Offset         int    `json:"offset,omitempty"`
}