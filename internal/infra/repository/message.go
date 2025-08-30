package repository

import (
	"context"
	"time"
	"fmt"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"	
	"github.com/jackc/pgx/v5/pgxpool"
	
)

type messageRepository struct {
    db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *messageRepository {
    return &messageRepository{db: db}
}
func (r *messageRepository) SaveMessage(ctx context.Context, message dto.Message) (*dto.Message, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    query := `
        INSERT INTO messages (
            conversation_id, sender_id, receiver_id, sender_type, 
            content, attachment_url, message_type, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at
    `
    
    row := r.db.QueryRow(
        ctx, query,
        message.ConversationID, message.SenderID, message.ReceiverID, message.SenderType,
        message.Content, message.AttachmentURL, message.MessageType.String(), time.Now(),
    )
    
    err := row.Scan(&message.ID, &message.CreatedAt)
    if err != nil {
        return nil, err
    }
    
    return &message, nil
}
func (r *messageRepository) GetMessagesByConversation(ctx context.Context, conversationID string, limit, offset int) ([]dto.Message, error) {
	 ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    query := `
        SELECT id, conversation_id, sender_id, receiver_id, sender_type, 
               content, attachment_url, message_type, created_at, read_at, is_deleted
        FROM messages
        WHERE conversation_id = $1 AND is_deleted = false
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `
    
    rows, err := r.db.Query(ctx, query, conversationID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var messages []dto.Message
    for rows.Next() {
        var msg dto.Message
        var msgType string
        
        err := rows.Scan(
            &msg.ID, &msg.ConversationID, &msg.SenderID, &msg.ReceiverID, &msg.SenderType,
            &msg.Content, &msg.AttachmentURL, &msgType, &msg.CreatedAt, &msg.ReadAt, &msg.IsDeleted,
        )
        if err != nil {
            return nil, err
        }
        
        msg.MessageType = dto.MessageType(msgType)
        messages = append(messages, msg)
    }
    
    if err := rows.Err(); err != nil {
        return nil, err
    }
    
    return messages, nil
}
func (r *messageRepository) GetUnreadMessageCount(ctx context.Context, userID int64, isAdmin bool) (int, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    query := `
        SELECT COUNT(*) 
        FROM messages 
        WHERE receiver_id = $1 
        AND read_at IS NULL 
        AND is_deleted = false
    `
    
    if isAdmin {
        query += ` AND sender_type = 'user'`
    } else {
        query += ` AND sender_type = 'admin'`
    }

    var count int
    err := r.db.QueryRow(ctx, query, userID).Scan(&count)
    if err != nil {
        return 0, err
    }

    return count, nil
}

func (r *messageRepository) MarkAsRead(ctx context.Context, messageID int64) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    query := `
        UPDATE messages 
        SET read_at = NOW() 
        WHERE id = $1 AND read_at IS NULL
    `

    result, err := r.db.Exec(ctx, query, messageID)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("message with ID %d not found or already read", messageID)
    }

    return nil
}

func (r *messageRepository) GetRecentConversations(ctx context.Context, userID int64, isAdmin bool, limit, offset int) ([]dto.Message, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    query := `
        WITH LastMessages AS (
            SELECT DISTINCT ON (conversation_id) 
                id, conversation_id, sender_id, receiver_id, sender_type,
                content, attachment_url, message_type, created_at, read_at, is_deleted
            FROM messages
            WHERE (sender_id = $1 OR receiver_id = $1)
            AND is_deleted = false
    `

    if isAdmin {
        query += ` AND (sender_type = 'user' OR (sender_type = 'admin' AND sender_id = $1))`
    } else {
        query += ` AND (sender_type = 'admin' OR (sender_type = 'user' AND sender_id = $1))`
    }

    query += `
            ORDER BY conversation_id, created_at DESC
        )
        SELECT * FROM LastMessages
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

    rows, err := r.db.Query(ctx, query, userID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var conversations []dto.Message
    for rows.Next() {
        var msg dto.Message
        var msgType string

        err := rows.Scan(
            &msg.ID, &msg.ConversationID, &msg.SenderID, &msg.ReceiverID, &msg.SenderType,
            &msg.Content, &msg.AttachmentURL, &msgType, &msg.CreatedAt, &msg.ReadAt, &msg.IsDeleted,
        )
        if err != nil {
            return nil, err
        }

        msg.MessageType = dto.MessageType(msgType)
        conversations = append(conversations, msg)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return conversations, nil
}