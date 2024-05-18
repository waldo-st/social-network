package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"social/internal/core/domain"
	"social/internal/core/port"
	util "social/internal/core/utils"
	"time"

	"github.com/gorilla/websocket"
)

type chatService struct {
	repo    port.ChatRepository
	timeout time.Duration
}

func NewChatServive(r port.ChatRepository) *chatService {
	return &chatService{r, time.Duration(3) * time.Second}
}

func (s *chatService) CreateMessage(ctx context.Context, message *domain.Chat) (*domain.Chat, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)

	// validate the sender and group id's
	defer cancel()

	return s.repo.CreateMessage(ctx, message)
}

func (s *chatService) GetChatsByGroupId(ctx context.Context, id int) ([]*domain.Chat, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.repo.GetChatsByGroupId(ctx, id)
}

func (s *chatService) ReadMessage(ctx context.Context, client *domain.Client, hub *domain.Hub) {

	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var wsMsg *domain.WsChatReq
		err = json.Unmarshal(p, &wsMsg)
		if err != nil {
			log.Println("JSON Unmarshal:", err)
			continue
		}
		chat := &domain.Chat{
			SenderId:  client.Id,
			Username:  client.Username,
			GroupId:   client.RoomId,
			Content:   wsMsg.Content,
			Type:      wsMsg.Type,
			CreatedAt: time.Now(),
		}
		fmt.Println("Chat => ", chat)
		if _, err = s.repo.CreateMessage(context.Background(), chat); err != nil {
			log.Println("Error saving msg:", err)
			return
		}

		hub.Broadcast <- chat
	}

}

func (s *chatService) WriteMessage(client *domain.Client) {
	defer func() {
		client.Conn.Close()
	}()

	for {
		msg, ok := <-client.Message
		if !ok {
			return
		}
		// fmt.Println("Client => ", client.Conn, "Message => ", msg)
		util.Broadcast(client.Conn, msg)
	}
}
