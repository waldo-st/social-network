package handler

import (
	"fmt"
	"net/http"
	"social/app"
	"social/internal/core/domain"
	"social/internal/core/port"
	"strconv"
)

type FollowHandler struct {
	srv port.FollowerService
	nh  NotificationHandler
}

var Message = map[string]string{
	"follow":     "want to follow you",
	"joingroup":  "want to join the group",
	"invitation": "invite you to join the groupe",
	"post":       "post a new post",
	"event":      "post a new event",
}

func NewFollowerHandler(srv port.FollowerService, nsrv NotificationHandler) *FollowHandler {
	return &FollowHandler{srv, nsrv}
}

func (h *FollowHandler) FollowRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")

		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		sender, ok := value.(int)
		if !ok {
			ctx.HandleError("invalid sender id", http.StatusInternalServerError)
			return
		}

		receiver, err := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
		if err != nil {
			ctx.HandleError("invalid receiver id parameter", http.StatusBadRequest)
			return
		}

		u, err := h.nh.usrv.GetUserById(ctx.Request.Context(), receiver)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		if sender == receiver {
			ctx.HandleError("the sender and the receiver can't be the same user", http.StatusBadRequest)
			return
		}

		typ := ctx.Request.URL.Query().Get("type")

		if u.IsPublic {
			f := &domain.Follow{
				FollowerId: sender,
				FolloweeId: receiver,
			}
			_, err := h.srv.Follow(ctx.Request.Context(), f)
			if err != nil {
				ctx.HandleError(err.Error(), http.StatusForbidden)
			}
			return
		}

		// check typ validation [françois]
		err = h.nh.Request(ctx, sender, receiver, 0, typ)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (h *FollowHandler) FollowReply() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}
		receiver, ok := value.(int)
		if !ok {
			ctx.HandleError("invalid sender id", http.StatusInternalServerError)
			return
		}

		nid, err := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
		if err != nil {
			ctx.HandleError("invalid id parameter", http.StatusBadRequest)
			return
		}

		status := ctx.Request.URL.Query().Get("status")

		n, err := h.nh.srv.Get(ctx.Request.Context(), nid)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("notification: ", n)
		if receiver != n.Receiver {
			ctx.HandleError("not allow", http.StatusForbidden)
			return
		}

		if status == "accepted" {
			f := &domain.Follow{
				FollowerId: n.Sender,
				FolloweeId: n.Receiver,
			}
			_, err := h.srv.Follow(ctx.Request.Context(), f)
			if err != nil {
				fmt.Println("errrrr", err)
				ctx.HandleError(err.Error(), http.StatusForbidden)
				return
			}
		}

		err = h.nh.Reply(ctx, n.Sender, nid, status)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusForbidden)
			return
		}
	}
}

func (h FollowHandler) ListFolloweeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		fmt.Println(value)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		userId, ok := value.(int)
		if !ok {
			ctx.HandleError("invalid id type", http.StatusInternalServerError)
			return
		}

		id, err := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
		if err != nil {
			ctx.HandleError("invalid receiver id parameter", http.StatusBadRequest)
			return
		}

		if id == 0 {
			id = userId
		}

		followee, err := h.srv.GetUserFollowee(ctx.Request.Context(), id)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("followee => ", followee)
		ctx.SendResponse(followee, http.StatusOK)
	}
}

func (h FollowHandler) ListFollowersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		fmt.Println(value)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		userId, ok := value.(int)
		if !ok {
			ctx.HandleError("invalid userId parameter", http.StatusInternalServerError)
			return
		}

		id, err := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
		if err != nil {
			ctx.HandleError("invalid receiver id parameter", http.StatusBadRequest)
			return
		}

		if id == 0 {
			id = userId
		}

		followers, err := h.srv.GetUserFollowers(ctx.Request.Context(), id)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		ctx.SendResponse(followers, http.StatusOK)
	}
}

func (h *FollowHandler) UnFollow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}
		follower, ok := value.(int)
		if !ok {
			ctx.HandleError("invalid sender id", http.StatusInternalServerError)
			return
		}

		followee, err := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		if follower == followee {
			ctx.HandleError("the follower and the followee can't be the same user", http.StatusBadRequest)
			return
		}

		if err := h.srv.UnFollow(ctx.Request.Context(), follower, followee); err != nil {
			ctx.HandleError(err.Error(), http.StatusForbidden)
			return
		}

		if err := h.nh.srv.Delete(ctx.Request.Context(), follower, followee); err != nil {
			ctx.HandleError(err.Error(), http.StatusUnauthorized)
			return
		}
	}
}
