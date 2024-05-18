package handler

import (
	"fmt"
	"net/http"
	"social/app"
	"social/internal/core/domain"
	"social/internal/core/port"
	"social/internal/core/service"
	util "social/internal/core/utils"
	"strconv"

	"github.com/gorilla/websocket"
)

type WsHandler struct {
	srv port.ChatService
	hub *domain.Hub
	ts  *service.TokenService
}

func NewWsHandler(h *domain.Hub, s port.ChatService, ts *service.TokenService) *WsHandler {
	return &WsHandler{hub: h, srv: s, ts: ts}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		// return origin == "http://localhost:3000"
		return true // for postmam test
	},
}

func (h *WsHandler) JoinRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)
		conn, err := upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		roomId, err := strconv.Atoi(ctx.Request.URL.Query().Get("roomId"))
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("ROOMID => ", roomId)
		t := ctx.Request.URL.Query().Get("token")
		if _, err := h.ts.VerifyToken(ctx.Request.Context(), t); err != nil {
			ctx.HandleError(err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := h.ts.GetUserByToken(ctx.Request.Context(), t)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		client := &domain.Client{
			Conn:     conn,
			Message:  make(chan *domain.Chat, 1024),
			Id:       token.UserId,
			RoomId:   roomId,
			Username: token.Username,
			Type:     "group",
		}

		fmt.Println("h", h.hub.Rooms)

		chats, err := h.srv.GetChatsByGroupId(ctx.Request.Context(), roomId)
		if err != nil {
			ctx.HandleError("Could not retrieve room chats", http.StatusInternalServerError)
			return
		}
		// h.hub.Register <- client

		util.Broadcast(client.Conn, chats)
		fmt.Println("Chat => ", chats)

		// go h.srv.ReadMessage(ctx.Request.Context(), client, h.hub)
		// go h.srv.WriteMessage(client)

	}
}

func (h *WsHandler) ChatHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		conn, err := upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		otherId, err := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		t := ctx.Request.URL.Query().Get("token")
		if _, err := h.ts.VerifyToken(ctx.Request.Context(), t); err != nil {
			ctx.HandleError(err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := h.ts.GetUserByToken(ctx.Request.Context(), t)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		roomId := token.UserId + otherId

		if _, ok := h.hub.PrivateRooms[roomId]; !ok {
			h.hub.PrivateRooms[roomId] = &domain.PrivateRoom{
				Id:      roomId,
				Name:    fmt.Sprintf("%d", roomId),
				Clients: make(map[int]*domain.Client),
			}
		}

		client := &domain.Client{
			Conn:     conn,
			Message:  make(chan *domain.Chat, 1024),
			Id:       token.UserId,
			RoomId:   roomId,
			Username: token.Username,
			Type:     "private",
		}

		h.hub.Register <- client

		chats, err := h.srv.GetChatsByGroupId(ctx.Request.Context(), roomId)
		if err != nil {
			ctx.HandleError("Could not retrieve room chats", http.StatusInternalServerError)
			return
		}
		fmt.Println("ChatPrivate => ", chats)
		go h.srv.ReadMessage(ctx.Request.Context(), client, h.hub)
		go h.srv.WriteMessage(client)

		util.Broadcast(client.Conn, chats)
	}
}
