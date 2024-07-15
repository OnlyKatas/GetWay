package run

import (
	"context"
	"getway/config"
	"getway/internal/modules/controller"
	"getway/internal/modules/grpc/auth"
	payGRPC "getway/internal/modules/grpc/pay"
	"getway/internal/modules/grpc/user"
	"getway/internal/router"
	"log"
	"net/http"
	_ "net/http/pprof"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(cfg *config.Config) error {
	ctx := context.Background()
	authClient, err := auth.NewAuthClient(ctx, cfg.GRPC.AuthAddress)
	if err != nil {
		log.Fatalf("auth.NewAuthClient failed: %v", err)
	}
	userClient, err := user.NewUserClient(ctx, cfg.GRPC.UserAddress)
	if err != nil {
		log.Fatalf("new pay client failed: %v", err)
	}
	payClient, err := payGRPC.NewClient(ctx, cfg.GRPC.PayAddress)
	if err != nil {
		log.Fatalf("new pay client failed: %v", err)
	}

	handlers := controller.NewHandler(authClient, userClient, payClient)
	rout := new(router.Router)
	s.httpServer = &http.Server{
		Addr:    cfg.Http.Port,
		Handler: rout.InitRoutes(handlers),
	}
	log.Printf("Listening on port %s", cfg.Http.Port)
	return s.httpServer.ListenAndServe()
}

func (a *Server) Stop() {
	const op = "grpcapp.Stop"
	log.Println("Stopping " + op)

	a.httpServer.Shutdown(context.Background())
}
