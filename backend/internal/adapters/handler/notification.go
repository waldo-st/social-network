package handler

import (
	"fmt"
	"net/http"
	"social/app"
	"social/internal/core/domain"
	"social/internal/core/port"
	util "social/internal/core/utils"
	"strings"
)

type NotificationHandler struct {
	srv  port.NotificationService
	usrv port.UserService
}

func NewNotificationHandler(srv port.NotificationService, usrv port.UserService) *NotificationHandler {
	return &NotificationHandler{srv, usrv}
}

var Clients = make(map[int]domain.Client, 0)

func (n *NotificationHandler) WsHandshake() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		token := ctx.Request.URL.Query().Get("token")
		if strings.TrimSpace(token) == "" {
			ctx.HandleError("unauthorized", http.StatusUnauthorized)
			return
		}

		userId, ok := Sessions[token]
		if !ok {
			ctx.HandleError("unauthorized", http.StatusUnauthorized)
			return
		}

		ws, err := upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
		if err != nil {
			ctx.HandleError("can not init connexion", http.StatusInternalServerError)
			return
		}

		Clients[userId] = domain.Client{
			Id:   userId,
			Conn: ws,
		}

		nts, err := n.srv.GetNotifications(ctx.Request.Context(), userId)
		if err != nil {
			ctx.HandleError("can not init connexion", http.StatusInternalServerError)
			return
		}

		util.Broadcast(ws, nts)
		fmt.Println("NTS => ", nts)
	}
}

var message = map[string]string{
	"follow": " want to follow you ",
	"join":   " want to join the group ",
	"invite": " invite you to join the group ",
	"post":   " post a new post ",
	"event":  " post a new event ",
}

func (h *NotificationHandler) Request(ctx *app.Context, sender, receiver, gid int, typ string) error {
	n := &domain.Notification{
		GroupId:  gid,
		Sender:   sender,
		Receiver: receiver,
		Status:   "pending",
		Type:     typ,
		Message:  message[typ],
	}

	_, err := h.srv.Add(ctx.Request.Context(), n)
	if err != nil {
		ctx.HandleError(err.Error(), http.StatusInternalServerError)
		return err
	}

	notifs, err := h.srv.GetNotifications(ctx.Request.Context(), receiver)
	if err != nil {
		ctx.HandleError(err.Error(), http.StatusInternalServerError)
		return err
	}

	if client, ok := Clients[receiver]; ok {
		util.Broadcast(client.Conn, notifs)
		// err := Broadcast(client.Conn, notifs)
		// if err != nil {
		// 	delete(Clients, receiver)
		// 	return err
		// }
	}

	return nil
}

func (h *NotificationHandler) Reply(ctx *app.Context, userId, notifid int, status string) error {
	fmt.Println("hello")
	if err := h.srv.Update(ctx.Request.Context(), status, notifid); err != nil {
		ctx.HandleError(err.Error(), http.StatusBadRequest)
		return err
	}

	notifs, err := h.srv.GetNotifications(ctx.Request.Context(), userId)
	if err != nil {
		ctx.HandleError(err.Error(), http.StatusInternalServerError)
		return err
	}

	if client, ok := Clients[userId]; ok {
		util.Broadcast(client.Conn, notifs)
		// err := Broadcast(client.Conn, notifs)
		// if err != nil {
		// 	delete(Clients, userId)
		// 	return err
		// }

	}

	return nil
}
