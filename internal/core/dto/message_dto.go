package dto
import (
    "time"
)
type MessageType string

const (
    TextMessage    MessageType = "text"
    SystemMessage   MessageType = "system"
)

type Message struct {
    ID             int64                    `json:"id"`
    ConversationID string                   `json:"conversation_id"`
    SenderID       int64                    `json:"sender_id"`
    ReceiverID     int64                    `json:"receiver_id"`
    SenderType     string                   `json:"sender_type"` // "admin" hoặc "user"
    Content        string                   `json:"content"`
    AttachmentURL  string                   `json:"attachment_url,omitempty"`
    MessageType    MessageType               `json:"message_type"`
    CreatedAt      time.Time                `json:"created_at"`
    ReadAt         *time.Time               `json:"read_at,omitempty"`
    IsDeleted      bool                     `json:"is_deleted"`
}
func (mt MessageType) IsValid() bool {
    switch mt {
    case TextMessage,  SystemMessage:
        return true
    default:
        return false
    }
}

func (mt MessageType) String() string {
    return string(mt)
}