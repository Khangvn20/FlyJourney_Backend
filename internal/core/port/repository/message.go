package repository
import(
	"context"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
)
type MessageRepository interface {
	SaveMessage(ctx context.Context, message dto.Message) (*dto.Message, error)
	GetMessagesByConversation(ctx context.Context, conversationID string, limit, offset int) ([]dto.Message, error)
	GetUnreadMessageCount(ctx context.Context, userID int64, isAdmin bool) (int, error)
    MarkAsRead(ctx context.Context, messageID int64) error
	GetRecentConversations(ctx context.Context, userID int64, isAdmin bool, limit, offset int) ([]dto.Message, error)
}