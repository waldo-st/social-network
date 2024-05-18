package handler

import (
	"fmt"
	"io"
	"net/http"
	"social/app"
	"social/internal/core/domain"
	"social/internal/core/port"
	util "social/internal/core/utils"
	"strconv"
)

type GroupHandler struct {
	hub  domain.Hub
	gsrv port.GroupService
	ns   NotificationHandler
}

func NewGroupHandler(gs port.GroupService, h domain.Hub, ns NotificationHandler) *GroupHandler {
	return &GroupHandler{h, gs, ns}
}

func (g *GroupHandler) CreateGroupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)
		var body domain.Group

		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
		}

		id, ok := value.(int)
		if !ok {
			ctx.HandleError("invalid id type", http.StatusInternalServerError)
			return
		}

		if err := ctx.BindJson(&body); err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
		}

		body.Admin = id

		group, err := g.gsrv.CreateGroup(ctx.Request.Context(), &body)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		// register rom to hub
		g.hub.Rooms[group.Id] = &domain.Room{
			Id:      group.Id,
			Name:    group.Title,
			Clients: make(map[int]*domain.Client),
		}

		ctx.SendResponse(group, http.StatusAccepted)
	}
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func (g *GroupHandler) JoinRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		user, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
		}
		sender, ok := user.(int)
		if !ok {
			ctx.HandleError("invalid id type", http.StatusInternalServerError)
			return
		}

		gid, err := strconv.Atoi(ctx.Request.URL.Query().Get("gid"))
		if err != nil {
			ctx.HandleError("invalid receiver id parameter", http.StatusBadRequest)
			return
		}
		gr, err := g.gsrv.GetGroupById(ctx.Request.Context(), gid)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		typ := ctx.Request.URL.Query().Get("type")
		// b, _ := io.ReadAll(r.Body)
		var u []User
		err = ctx.BindJson(&u)
		if err != nil {
			if err != io.EOF {
				ctx.HandleError(err.Error(), http.StatusInternalServerError)
				return
			}
		}
		// si c'est demande d'ajout
		if sender == gr.Admin && typ == "join" { // && []user est vide
			ctx.HandleError("request : You are already member", http.StatusBadRequest)
			return
		}
		//fin si c'est demande d'ajout
		// n := &domain.Notification{
		// 	GroupId: gr.Id,
		// 	Sender:  sender,
		// 	// Receiver: gr.Admin, // a commenter
		// 	Status:  "pending",
		// 	Type:    "group",
		// 	Message: message["group"] + gr.Title,
		// }
		msg := message["invite"]
		if typ == "join" {
			msg = message["join"]
			u = []User{{Id: gr.Admin}}
		}

		// si []user est vide initialise le avec admin
		// for []user invite
		// if tableau invite n'est pas vide {
		// }
		for _, user := range u {
			fmt.Println("Body => ", user.Id)
			n := &domain.Notification{
				GroupId: gr.Id,
				Sender:  sender,
				// Receiver: gr.Admin, // a commenter
				Receiver: user.Id,
				Status:   "pending",
				Type:     typ,
				Message:  msg + gr.Title,
			}

			_, err = g.ns.srv.Add(ctx.Request.Context(), n)
			if err != nil {
				ctx.HandleError(err.Error(), http.StatusInternalServerError)
				return
			}

			// notifs, err := g.ns.srv.GetNotifications(ctx.Request.Context(), n.Receiver)
			notifs, err := g.ns.srv.GetNotifications(ctx.Request.Context(), user.Id)
			if err != nil {
				ctx.HandleError(err.Error(), http.StatusInternalServerError)
				return
			}

			// if client, ok := Clients[n.Receiver]; ok {
			if client, ok := Clients[user.Id]; ok {
				util.Broadcast(client.Conn, notifs)
				// err := Broadcast(client.Conn, notifs)
				// if err != nil {
				// 	delete(Clients, receiver)
				// 	return err
				// }
			}
			// fin for user invite
		}
	}
}

func (g *GroupHandler) JoinReply() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		user, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
		}
		receiver, ok := user.(int)
		if !ok {
			ctx.HandleError("invalid id type", http.StatusInternalServerError)
			return
		}

		id, err := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
		if err != nil {
			ctx.HandleError("invalid receiver id parameter", http.StatusBadRequest)
			return
		}

		status := ctx.Request.URL.Query().Get("status")
		typ := ctx.Request.URL.Query().Get("type")

		n, err := g.ns.srv.Get(ctx.Request.Context(), id)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		if n.Receiver != receiver {
			ctx.HandleError("user is not the admin", http.StatusForbidden)
			return
		}

		if status == "accepted" {
			var err error
			if typ == "invite" {
				err = g.gsrv.AddMember(ctx.Request.Context(), &domain.Group{Admin: n.Receiver, Id: n.GroupId})
			} else {
				err = g.gsrv.AddMember(ctx.Request.Context(), &domain.Group{Admin: n.Sender, Id: n.GroupId})
			}
			if err != nil {
				ctx.HandleError(err.Error(), http.StatusBadRequest)
				return
			}
		}

		if err := g.ns.srv.Update(ctx.Request.Context(), status, id); err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		notifs, err := g.ns.srv.GetNotifications(ctx.Request.Context(), n.Receiver)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		if client, ok := Clients[n.Receiver]; ok {
			util.Broadcast(client.Conn, notifs)
			// err := Broadcast(client.Conn, notifs)
			// if err != nil {
			// 	delete(Clients, userId)
			// 	return err
			// }
		}
	}
}

// func (g *GroupHandler) InviteRequest() http.HandlerFunc {

// }

type GrouMemberReq struct {
	Id int `json:"groupId"`
}

func (h *GroupHandler) ListGroupsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
		}
		id, ok := value.(int)
		if !ok {
			ctx.HandleError("invalid id type", http.StatusInternalServerError)
			return
		}

		fmt.Println("USER ID", id)

		groups, err := h.gsrv.ListGroups(ctx.Request.Context(), id)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		// add rooms for each group
		for _, group := range groups {
			if _, ok := h.hub.Rooms[group.Id]; !ok {
				h.hub.Rooms[group.Id] = &domain.Room{
					Id:      group.Id,
					Name:    group.Title,
					Clients: make(map[int]*domain.Client),
				}
			}
		}
		fmt.Println(h.hub.Rooms)
		ctx.SendResponse(groups, http.StatusAccepted)
	}
}

func (h *GroupHandler) JoinedGroupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
		}
		userId, ok := value.(int)
		if !ok {
			ctx.HandleError("invalid id type", http.StatusInternalServerError)
			return
		}

		groups, err := h.gsrv.ListJoinedGroups(ctx.Request.Context(), userId)
		if err != nil {
			ctx.HandleError("invalid id type", http.StatusInternalServerError)
			return
		}

		for _, group := range groups {
			if _, ok := h.hub.Rooms[group.Id]; !ok {
				h.hub.Rooms[group.Id] = &domain.Room{
					Id:      group.Id,
					Name:    group.Title,
					Clients: make(map[int]*domain.Client),
				}
			}
		}

		ctx.SendResponse(groups, http.StatusOK)
	}
}

func (h *GroupHandler) UnjoinedGroupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
		}

		userId, ok := value.(int)
		if !ok {
			ctx.HandleError("invalid id type", http.StatusInternalServerError)
			return
		}

		groups, err := h.gsrv.ListUnjoinedGroups(ctx.Request.Context(), userId)
		if err != nil {
			ctx.HandleError("invalid id type", http.StatusInternalServerError)
			return
		}

		ctx.SendResponse(groups, http.StatusOK)
	}
}
