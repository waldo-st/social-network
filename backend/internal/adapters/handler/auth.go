package handler

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"social/app"
// 	"social/internal/core/domain"
// 	"social/internal/core/port"
// )

// type AuthHandler struct {
// 	srv port.TokenService
// }

// func NewAuthHandler(srv port.TokenService) *AuthHandler {
// 	return &AuthHandler{srv}
// }

// func (h *AuthHandler) LoginHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := app.NewContext(w, r)

// 		var userLog domain.UserLog

// 		if err := ctx.BindJson(&userLog); err != nil {
// 			ctx.HandleError(err.Error(), http.StatusBadRequest)
// 		}

// 		u, err := h.srv.Login(ctx.Request.Context(), &userLog)
// 		if err != nil {
// 			ctx.HandleError(err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		token, err := h.srv.CreateToken(ctx.Request.Context(), u)
// 		if err != nil {
// 			fmt.Println("ERROR CREATING TOKEN:", err)
// 			ctx.HandleError(err.Error(), http.StatusUnauthorized)
// 			return
// 		}

// 		user := domain.UserRes{
// 			Id: u.Id,
// 			Email:       u.Email,
// 			FirstName:   u.FirstName,
// 			LastName:    u.LastName,
// 			DateOfBirth: u.DateOfBirth,
// 			Avatar:      u.Avatar,
// 			Nickname:    u.Nickname,
// 			About:       u.About,
// 			IsPublic:    u.IsPublic,
// 		}

// 		response, err := h.srv.GetResponseData(ctx.Request.Context(), u)
// 		if err != nil {
// 			log.Printf("ERROR: %v", err)
// 		}
// 		response.Profile = user
// 		response.Token = token

// 		ctx.SendResponse(response, http.StatusAccepted)
// 	}
// }

