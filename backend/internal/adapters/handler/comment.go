package handler

import (
	"fmt"
	"html"
	"net/http"
	"social/app"
	"social/internal/core/domain"
	"social/internal/core/port"
	"strconv"
	"time"
)

type CommentHandler struct {
	srv port.CommentService
}

func NewCommentHandler(s port.CommentService) *CommentHandler {
	return &CommentHandler{s}
}

func (ch *CommentHandler) CreateComment() http.HandlerFunc {
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

		postId, err := strconv.Atoi(ctx.Request.URL.Query().Get("postId"))
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		var commentReq domain.Comment
		if err := ctx.BindJson(&commentReq); err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
		}

		comment := domain.Comment{
			UserId:    userId,
			PostId:    postId,
			Content:   html.EscapeString(commentReq.Content),
			Image:     commentReq.Image,
			CreatedAt: time.Now(),
		}

		response, err := ch.srv.CreateComment(ctx.Request.Context(), &comment)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}
		ctx.SendResponse(response, http.StatusCreated)
	}
}

func (ch *CommentHandler) GetComments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		// value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		// if err != nil {
		// 	ctx.HandleError(err.Error(), http.StatusBadRequest)
		// }

		// _, ok := value.(int)
		// if !ok {
		// 	ctx.HandleError("invalid id type", http.StatusInternalServerError)
		// 	return
		// }

		postId, err := strconv.Atoi(ctx.Request.URL.Query().Get("postId"))
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Hello from comments", postId)

		commnents, err := ch.srv.GetCommentsByPostId(ctx.Request.Context(), postId)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("comments retrieved", commnents)
		ctx.SendResponse(commnents, http.StatusOK)
	}
}
