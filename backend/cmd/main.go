package main

import (
	"log"
	"os"
	"social/internal/adapters/handler"
	"social/internal/adapters/router"
	db "social/internal/adapters/storage/sqlite"
	"social/internal/adapters/storage/sqlite/repository"
	"social/internal/core/domain"
	"social/internal/core/service"
)

func main() {
	db, err := db.New()
	if err != nil {
		log.Println("Error initializing database connection", err)
		os.Exit(1)
	}

	defer db.Close()

	if err := db.Migrate(); err != nil {
		log.Println("Error migrating database", err)
		os.Exit(1)
	}

	var r = new(router.Route)

	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	r.UserHandler = handler.NewUserHandler(userService)

	// tokenRepo := repository.NewTokenRepository(db.DB)
	// r.TokenService = service.NewTokenService(userRepo, tokenRepo)

	notifRepo := repository.NewNotificationRepository(db.DB)
	notifService := service.NewNotificationService(notifRepo)
	r.NotifHandler = handler.NewNotificationHandler(notifService, userService)

	followerRepo := repository.NewFollowerRepository(db.DB)
	followerSrv := service.NewFollower(followerRepo)
	r.FollowHandler = handler.NewFollowerHandler(followerSrv, *r.NotifHandler)

	postRepo := repository.NewPostRepository(db.DB)
	postService := service.NewPostService(userRepo, postRepo)
	r.PostHandler = handler.NewPostHandler(postService)

	commentRepo := repository.NewCommentRepository(db.DB)
	comentService := service.NewCommentService(userRepo, postRepo, commentRepo)
	r.CommentHandler = handler.NewCommentHandler(comentService)

	hub := domain.NewHub()

	groupRepo := repository.NewGroupRepository(db.DB)
	groupService := service.NewGroupService(userRepo, groupRepo)
	r.GroupHandler = handler.NewGroupHandler(groupService, *hub, *r.NotifHandler)

	eventRepo := repository.NewEventRepository(db.DB)
	eventService := service.NewEventService(userRepo, groupRepo, eventRepo)
	r.EventHandler = handler.NewEventHandler(eventService)

	// clientService := service.NewClientService(chatRepository)
	chatRepository := repository.NewChatRepository(db.DB)
	chatService := service.NewChatServive(chatRepository)

	// r.WsHandler = handler.NewWsHandler()

	go hub.Run()

	tokenRepo := repository.NewTokenRepository(db.DB)
	tokenService := service.NewTokenService(userRepo, postRepo, commentRepo, groupRepo, followerRepo, eventRepo, tokenRepo)

	r.Authandler = handler.NewAuthHandler(tokenService)
	r.WsHandler = handler.NewWsHandler(hub, chatService, tokenService)
	router.InitRouter(r, tokenService)
}
