package handler

import (
	"net/http"
	"social/app"
	"social/internal/core/domain"
	"social/internal/core/port"
	"strconv"
	"time"
)

type EventHandler struct {
	srv port.EventService
}

func NewEventHandler(s port.EventService) *EventHandler {
	return &EventHandler{s}
}

func (h *EventHandler) CreateEventHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		var req domain.EventReq

		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		userId, ok := value.(int)
		if !ok {
			ctx.HandleError("invalid id type", http.StatusInternalServerError)
			return
		}

		if err := ctx.BindJson(&req); err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		groupId, err := strconv.Atoi(ctx.Request.URL.Query().Get("groupId"))
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		event := &domain.Event{
			GroupId:     groupId,
			CreatorId:   userId,
			Title:       req.Title,
			Description: req.Description,
			CreatedAt:   time.Now(),
		}

		response, err := h.srv.CreateEvent(ctx.Request.Context(), event)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}
		ctx.SendResponse(response, http.StatusOK)
	}
}

func (h *EventHandler) ListGroupEventsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)
		groupId, err := strconv.Atoi(ctx.Request.URL.Query().Get("groupId"))
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return

		}

		events, err := h.srv.ListGroupEvents(ctx.Request.Context(), groupId)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}
		ctx.SendResponse(events, http.StatusOK)
	}
}

func (h *EventHandler) RespondToEventHandler() http.HandlerFunc {
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

		var req *domain.EventReaction

		if err := ctx.BindJson(&req); err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		reaction := &domain.Reaction{
			EventId:   req.Id,
			UserId:    userId,
			Status:    req.Status,
			CreatedAt: time.Now(),
		}

		response, err := h.srv.RespondToEvent(ctx.Request.Context(), reaction)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}
		ctx.SendResponse(response, http.StatusOK)
	}
}
