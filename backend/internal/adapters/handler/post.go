package handler

import (
	"fmt"
	"net/http"
	"social/app"
	"social/internal/core/domain"
	"social/internal/core/port"
	"strconv"
	"time"
)

type PostHandler struct {
	srv port.PostService
}

func NewPostHandler(s port.PostService) *PostHandler {
	return &PostHandler{s}
}

func (ph *PostHandler) CreatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		var post domain.Post

		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
		}

		id, ok := value.(int)
		if !ok {
			ctx.HandleError("invalid id type", http.StatusInternalServerError)
			return
		}

		// groupId, err := strconv.Atoi(ctx.Request.URL.Query().Get("postId"))
		// if err != nil {
		// 	ctx.HandleError("invalid groupId type", http.StatusInternalServerError)
		// 	return
		// }
		groupId, err := strconv.Atoi(ctx.Request.URL.Query().Get("groupId"))
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		if err := ctx.BindJson(&post); err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
			return
		}

		post.UserId = id
		post.GroupId = groupId
		post.CreatedAt = time.Now()

		_, err = ph.srv.CreatePost(ctx.Request.Context(), &post)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}
		ctx.SendResponse(post, http.StatusOK)
	}
}

func (ph *PostHandler) ListPostsHandler() http.HandlerFunc {
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

		posts, err := ph.srv.ListPosts(ctx.Request.Context(), id)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("HEllo from posts", posts)
		ctx.SendResponse(posts, http.StatusOK)
	}
}

func (ph *PostHandler) GetUserPosts() http.HandlerFunc {
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

		posts, err := ph.srv.GetsPostsByUserId(ctx.Request.Context(), userId)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}
		ctx.SendResponse(posts, http.StatusOK)
	}
}

func (ph *PostHandler) GetGroupPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		groupId, err := strconv.Atoi(ctx.Request.URL.Query().Get("groupId"))
		if err != nil {
			ctx.HandleError("invalid id type", http.StatusInternalServerError)
			return
		}
		fmt.Println("Groupe => ", groupId)
		posts, err := ph.srv.GetPostsByGroupId(ctx.Request.Context(), groupId)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}
		ctx.SendResponse(posts, http.StatusOK)
	}
}
