package service

/*import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/repository"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/common/utils"
)
type chatService struct {
	messageRepo repository.MessageRepository
	redisService service.RedisService
}

func NewChatService(messageRepo repository.MessageRepository, redisService service.RedisService) service.ChatService {
	return &chatService{
		messageRepo: messageRepo,
		redisService: redisService,
	}
}
func (s *chatService) SendMessage(ctx context.Context, req *request.SendMessageRequest) *response.Response {
    // Validate input
    if req.SenderID <= 0 {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid sender ID",
        }
    }
            if req.ReceiverID <= 0 {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid receiver ID",
        }
    }

    if req.Content == "" {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Message content cannot be empty",
        }
    }

    // Validate message type
    messageType := dto.MessageType(req.MessageType)
    if !messageType.IsValid() {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid message type",
        }
    }

    // Generate conversation ID if not provided
    conversationID := req.ConversationID
    if conversationID == "" {
        conversationID = utils.GenerateConversationID(req.SenderID, req.ReceiverID)
    }

    // Create message DTO
    message := dto.Message{
        ConversationID: conversationID,
        SenderID:       req.SenderID,
        ReceiverID:     req.ReceiverID,
        SenderType:     req.SenderType,
        Content:        req.Content,
        AttachmentURL:  req.AttachmentURL,
        MessageType:    messageType,
        CreatedAt:      time.Now(),
        IsDeleted:      false,
    }

    // Save message to database
    savedMessage, err := s.messageRepo.SaveMessage(ctx, message)
    if err != nil {
        log.Printf("Error saving message: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to save message",
        }
    }

    // Cache message in Redis for real-time access
    if err := s.cacheMessage(ctx, savedMessage); err != nil {
        log.Printf("Error caching message: %v", err)
        // Don't return error for caching failure, message is already saved in DB
    }

    // Cache conversation list for both sender and receiver
    s.invalidateConversationCache(req.SenderID, req.ReceiverID)

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Message sent successfully",
        Data:         savedMessage,
    }
}
func (s *chatService) GetMessagesByConversation(ctx context.Context, req *request.GetMessagesRequest) *response.Response {
    if req.ConversationID == "" {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Conversation ID is required",
        }
    }

    // Set default pagination
    if req.Limit <= 0 {
        req.Limit = 20
    }
    if req.Offset < 0 {
        req.Offset = 0
    }

    // Try to get from cache first
    cacheKey := fmt.Sprintf("conversation:%s:page:%d:%d", req.ConversationID, req.Limit, req.Offset)
    if cachedMessages, err := s.getCachedMessages(ctx, cacheKey); err == nil {
        return &response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            ErrorMessage: "Messages retrieved successfully (from cache)",
            Data:         cachedMessages,
        }
    }

    // Get from database
    messages, err := s.messageRepo.GetMessagesByConversation(ctx, req.ConversationID, req.Limit, req.Offset)
    if err != nil {
        log.Printf("Error getting messages: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to retrieve messages",
        }
    }

    // Cache the result
    if err := s.cacheMessages(ctx, cacheKey, messages); err != nil {
        log.Printf("Error caching messages: %v", err)
    }

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Messages retrieved successfully",
        Data:         messages,
    }
}

// GetRecentConversations lấy danh sách conversation gần đây
func (s *chatService) GetRecentConversations(ctx context.Context, req *request.GetConversationsRequest) *response.Response {
    if req.UserID <= 0 {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid user ID",
        }
    }

    // Set default pagination
    if req.Limit <= 0 {
        req.Limit = 10
    }
    if req.Offset < 0 {
        req.Offset = 0
    }

    // Try to get from cache first
    cacheKey := fmt.Sprintf("conversations:user:%d:admin:%t:page:%d:%d", req.UserID, req.IsAdmin, req.Limit, req.Offset)
    if cachedConversations, err := s.getCachedMessages(ctx, cacheKey); err == nil {
        return &response.Response{
            Status:       true,
            ErrorCode:    error_code.Success,
            ErrorMessage: "Conversations retrieved successfully (from cache)",
            Data:         cachedConversations,
        }
    }

    // Get from database
    conversations, err := s.messageRepo.GetRecentConversations(ctx, req.UserID, req.IsAdmin, req.Limit, req.Offset)
    if err != nil {
        log.Printf("Error getting conversations: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to retrieve conversations",
        }
    }

    // Cache the result
    if err := s.cacheMessages(ctx, cacheKey, conversations); err != nil {
        log.Printf("Error caching conversations: %v", err)
    }

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Conversations retrieved successfully",
        Data:         conversations,
    }
}

// MarkAsRead đánh dấu tin nhắn đã đọc
func (s *chatService) MarkAsRead(ctx context.Context, req *request.MarkAsReadRequest) *response.Response {
    if req.MessageID <= 0 {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid message ID",
        }
    }

    err := s.messageRepo.MarkAsRead(ctx, req.MessageID)
    if err != nil {
        log.Printf("Error marking message as read: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to mark message as read",
        }
    }

    // Invalidate cache for unread count
    if req.UserID > 0 {
        cacheKey := fmt.Sprintf("unread_count:user:%d", req.UserID)
        s.redisService.Del(cacheKey)
    }

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Message marked as read successfully",
    }
}

func (s *chatService) GetUnreadCount(ctx context.Context, req *request.GetUnreadCountRequest) *response.Response {
    if req.UserID <= 0 {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InvalidRequest,
            ErrorMessage: "Invalid user ID",
        }
    }

    // Try to get from cache first
    cacheKey := fmt.Sprintf("unread_count:user:%d:admin:%t", req.UserID, req.IsAdmin)
    if cachedCount, err := s.redisService.Get(cacheKey); err == nil {
        var count int
        if json.Unmarshal([]byte(cachedCount), &count) == nil {
            return &response.Response{
                Status:       true,
                ErrorCode:    error_code.Success,
                ErrorMessage: "Unread count retrieved successfully (from cache)",
                Data:         map[string]interface{}{"unread_count": count},
            }
        }
    }

    // Get from database
    count, err := s.messageRepo.GetUnreadMessageCount(ctx, req.UserID, req.IsAdmin)
    if err != nil {
        log.Printf("Error getting unread message count: %v", err)
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Failed to get unread message count",
        }
    }

    // Cache the result for 5 minutes
    if countJSON, err := json.Marshal(count); err == nil {
        s.redisService.Set(cacheKey, string(countJSON), 5*time.Minute)
    }

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Unread count retrieved successfully",
        Data:         map[string]interface{}{"unread_count": count},
    }
}
// Helper methods for caching

func (s *chatService) cacheMessage(ctx context.Context, message *dto.Message) error {
    messageJSON, err := json.Marshal(message)
    if err != nil {
        return err
    }

    // Cache individual message
    messageKey := fmt.Sprintf("message:%d", message.ID)
    if err := s.redisService.Set(messageKey, string(messageJSON), 24*time.Hour); err != nil {
        return err
    }

    // Add to conversation stream
    conversationKey := fmt.Sprintf("conversation_stream:%s", message.ConversationID)
    if err := s.redisService.LPush(conversationKey, string(messageJSON)); err != nil {
        return err
    }

    // Keep only last 100 messages in stream
    s.redisService.LTrim(conversationKey, 0, 99)

    return nil
}

func (s *chatService) getCachedMessages(ctx context.Context, cacheKey string) ([]dto.Message, error) {
    cachedData, err := s.redisService.Get(cacheKey)
    if err != nil {
        return nil, err
    }

    var messages []dto.Message
    if err := json.Unmarshal([]byte(cachedData), &messages); err != nil {
        return nil, err
    }

    return messages, nil
}

func (s *chatService) cacheMessages(ctx context.Context, cacheKey string, messages []dto.Message) error {
    messagesJSON, err := json.Marshal(messages)
    if err != nil {
        return err
    }

    return s.redisService.Set(cacheKey, string(messagesJSON), 10*time.Minute)
}

func (s *chatService) invalidateConversationCache(senderID, receiverID int64) {
    // Invalidate conversation cache for both users
    patterns := []string{
        fmt.Sprintf("conversations:user:%d:*", senderID),
        fmt.Sprintf("conversations:user:%d:*", receiverID),
    }

    for _, pattern := range patterns {
        keys, err := s.redisService.Keys(pattern)
        if err == nil {
            for _, key := range keys {
                s.redisService.Del(key)
            }
        }
    }
}*/