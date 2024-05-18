package handler

import (
	"net/http"
	"social/app"
	"social/internal/core/domain"
	"social/internal/core/port"
	"strconv"
)

type UserHandler struct {
	usrv port.UserService
}

type AuthHandler struct {
	srv port.TokenService
}

// var Clients = make(map[int]domain.Client, 0)
// var Sessions = make(map[string]int)

func NewAuthHandler(srv port.TokenService) *AuthHandler {
	return &AuthHandler{srv}
}

func NewUserHandler(us port.UserService) *UserHandler {
	return &UserHandler{us}
}

func (h *UserHandler) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		var user domain.User
		user.IsPublic = true
		if err := ctx.BindJson(&user); err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
		}

		err := h.usrv.Register(ctx.Request.Context(), &user)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		ctx.SendResponse(map[string]interface{}{"message": "user created successfully", "status": http.StatusCreated}, http.StatusCreated)
	}

}

func (h *UserHandler) UpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		var user domain.User

		if err := ctx.BindJson(&user); err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
		}

		u, err := h.usrv.UpdateUser(ctx.Request.Context(), &user)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		ctx.SendResponse(u, http.StatusAccepted)
	}
}

var Sessions = make(map[string]int)

func (h *AuthHandler) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := app.NewContext(w, r)

		var userLog domain.UserLog

		if err := ctx.BindJson(&userLog); err != nil {
			ctx.HandleError(err.Error(), http.StatusBadRequest)
		}

		u, err := h.srv.Login(ctx.Request.Context(), &userLog)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := h.srv.CreateToken(ctx.Request.Context(), u)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusUnauthorized)
			return
		}

		Sessions[token] = u.Id

		ctx.SendResponse(token, http.StatusAccepted)
	}
}

// func (h *UserHandler) GetProfilHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := app.NewContext(w, r)
// 		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
// 		if err != nil {
// 			ctx.HandleError(err.Error(), http.StatusBadRequest)
// 		}
// 		userId, ok := value.(int)
// 		if !ok {
// 			ctx.HandleError("invalid id type", http.StatusInternalServerError)
// 			return
// 		}
// 		user, err := h.usrv.GetUserById(ctx.Request.Context(), userId)
// 		if err != nil {
// 			ctx.HandleError(err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		posts, err := h.usrv.GetOwnPosts(ctx.Request.Context(), userId)
// 		if err != nil {
// 			ctx.HandleError(err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		response := struct {
// 			Profil *domain.UserRes
// 			Posts  []*domain.PostInfo
// 		}{
// 			Profil: user,
// 			Posts:  posts,
// 		}
// 		ctx.SendResponse(response, http.StatusOK)
// 	}
// }

func (h *UserHandler) GetProfilHandler() http.HandlerFunc {
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

		id, err := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		if id == 0 {
			id = userId
		}

		user, err := h.usrv.GetUserById(ctx.Request.Context(), id)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		ok, err = h.usrv.CanSeeProfil(ctx.Request.Context(), id, user)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		if !ok {
			ctx.HandleError("please follow first", http.StatusUnauthorized)
			return
		}

		posts, err := h.usrv.GetOwnPosts(ctx.Request.Context(), id)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		response := struct {
			Profil *domain.UserRes
			Posts  []*domain.PostInfo
		}{
			Profil: user,
			Posts:  posts,
		}

		ctx.SendResponse(response, http.StatusOK)
	}
}

func (h *UserHandler) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *UserHandler) HomeHandler() http.HandlerFunc {
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

		user, err := h.usrv.GetUserById(ctx.Request.Context(), userId)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		users, err := h.usrv.ListUsers(ctx.Request.Context(), userId)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		events, err := h.usrv.ListAttendingEvents(ctx.Request.Context(), userId)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		connections, err := h.usrv.ListConnections(ctx.Request.Context(), userId)
		if err != nil {
			ctx.HandleError(err.Error(), http.StatusInternalServerError)
			return
		}

		response := struct {
			Profil      *domain.UserRes
			Users       []*domain.UserRes
			Connections []*domain.UserRes
			Events      []*domain.Event
		}{
			Profil:      user,
			Users:       users,
			Connections: connections,
			Events:      events,
		}

		ctx.SendResponse(response, http.StatusOK)
	}
}

// func (h *UserHandler) JoinedGroupHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := app.NewContext(w, r)
// 		value, err := ctx.GetContextValue(ctx.Request.Context(), "userId")
// 		if err != nil {
// 			ctx.HandleError(err.Error(), http.StatusBadRequest)
// 		}
// 		userId, ok := value.(int)
// 		if !ok {
// 			ctx.HandleError("invalid id type", http.StatusInternalServerError)
// 			return
// 		}
// 		groups, err := h.usrv.ListJoinedGroups(ctx.Request.Context(), userId)
// 		if err != nil {
// 			ctx.HandleError("invalid id type", http.StatusInternalServerError)
// 			return
// 		}
// 		ctx.SendResponse(groups, http.StatusOK)
// 	}
// }
