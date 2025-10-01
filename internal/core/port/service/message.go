package service
import (
	"context"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
)
type ChatService interface {
    SendMessage(ctx context.Context, req *request.SendMessageRequest) *response.Response
    GetConversationHistory(ctx context.Context, req *request.GetConversationHistoryRequest) *response.Response
    GetUserConversations(ctx context.Context, req *request.GetConversationsRequest) *response.Response
    MarkMessagesAsRead(ctx context.Context, req *request.MarkAsReadRequest) *response.Response
    GetUnreadCount(ctx context.Context, req *request.GetUnreadCountRequest) *response.Response
}