package router

import (
	"net/http"
	"social/app"
	"social/internal/adapters/handler"
	"social/internal/adapters/middleware"
	"social/internal/core/service"
)

type Route struct {
	UserHandler    *handler.UserHandler
	FollowHandler  *handler.FollowHandler
	PostHandler    *handler.PostHandler
	CommentHandler *handler.CommentHandler
	GroupHandler   *handler.GroupHandler
	Authandler     *handler.AuthHandler
	WsHandler      *handler.WsHandler
	EventHandler   *handler.EventHandler
	NotifHandler   *handler.NotificationHandler
	// TokenService   *service.TokenService
}

func InitRouter(r *Route, ts *service.TokenService) {
	rt := app.New(*ts)
	rt.HandleFunc("POST", "/register", r.UserHandler.RegisterHandler())
	rt.HandleFunc("POST", "/login", r.Authandler.LoginHandler())
	rt.HandleFunc("GET", "/profil", middleware.Auth(ts, r.UserHandler.GetProfilHandler()))
	// rt.HandleFunc("GET", "/followee/profil", middleware.Auth(ts, r.UserHandler.GetFolloweeProfilHandler()))

	rt.HandleFunc("GET", "/home", middleware.Auth(ts, r.UserHandler.HomeHandler()))

	rt.HandleFunc("POST", "/follow/request", middleware.Auth(ts, r.FollowHandler.FollowRequest()))
	rt.HandleFunc("POST", "/follow/reply", middleware.Auth(ts, r.FollowHandler.FollowReply()))
	rt.HandleFunc("GET", "/followers", middleware.Auth(ts, r.FollowHandler.ListFollowersHandler()))
	rt.HandleFunc("GET", "/followee", middleware.Auth(ts, r.FollowHandler.ListFolloweeHandler()))
	rt.HandleFunc("POST", "/unfollow", middleware.Auth(ts, r.FollowHandler.UnFollow()))

	rt.HandleFunc("GET", "/post", middleware.Auth(ts, r.PostHandler.ListPostsHandler()))
	rt.HandleFunc("POST", "/post", middleware.Auth(ts, r.PostHandler.CreatePostHandler()))
	rt.HandleFunc("POST", "/comment", middleware.Auth(ts, r.CommentHandler.CreateComment()))
	rt.HandleFunc("GET", "/comment", middleware.Auth(ts, r.CommentHandler.GetComments()))

	rt.HandleFunc("GET", "/group/posts", middleware.Auth(ts, r.PostHandler.GetGroupPostHandler()))

	// rt.HandleFunc("GET", "/group", middleware.Auth(ts, r.GroupHandler.ListGroupsHandler()))
	rt.HandleFunc("GET", "/joinedGroups", middleware.Auth(ts, r.GroupHandler.JoinedGroupHandler()))
	rt.HandleFunc("GET", "/group", middleware.Auth(ts, r.GroupHandler.ListGroupsHandler()))
	rt.HandleFunc("GET", "/unjoinedGroups", middleware.Auth(ts, r.GroupHandler.UnjoinedGroupHandler()))
	rt.HandleFunc("POST", "/group", middleware.Auth(ts, r.GroupHandler.CreateGroupHandler()))
	rt.HandleFunc("POST", "/joinGroup/request", middleware.Auth(ts, r.GroupHandler.JoinRequest()))
	rt.HandleFunc("POST", "/joinGroup/reply", middleware.Auth(ts, r.GroupHandler.JoinReply()))
	// rt.HandleFunc("GET", "/invite", middleware.Auth(ts, r.GroupHandler.Invite()))

	// rt.HandleFunc("GET", "/user/posts/", middleware.Auth(ts, r.PostHandler.GetUserPostsHandler()))

	// rt.HandleFunc("GET", "/users", middleware.Auth(ts, r.UserHandler.ListUsersHandler()))
	rt.HandleFunc("POST", "/event", middleware.Auth(ts, r.EventHandler.CreateEventHandler()))
	rt.HandleFunc("GET", "/event", middleware.Auth(ts, r.EventHandler.ListGroupEventsHandler()))
	rt.HandleFunc("POST", "/event/response", middleware.Auth(ts, r.EventHandler.RespondToEventHandler()))
	// rt.HandleFunc("POST", "/ws/ckreateRoom", r.WsHandler.CreateRoom())
	rt.HandleFunc("GET", "/groupChat", r.WsHandler.JoinRoomHandler())
	rt.HandleFunc("GET", "/chat", r.WsHandler.ChatHandler())

	// rt.HandleFunc("GET", "/ws/listRooms", r.WsHandler.ListRoomsHandler())
	// rt.HandleFunc("GET", "/ws/listClients", r.WsHandler.GetClients())

	rt.HandleFunc("GET", "/handshake", r.NotifHandler.WsHandshake())

	http.ListenAndServe(":8080", middleware.CORS(rt))
}
