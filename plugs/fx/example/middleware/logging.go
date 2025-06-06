package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/fitan/genx/plugs/fx/example/services"
)

type loggingUserService struct {
	next   services.UserService
	logger *slog.Logger
}

// @fx decorate services.UserService
func NewLoggingUserService(svc services.UserService, logger *slog.Logger) services.UserService {
	return &loggingUserService{
		next:   svc,
		logger: logger,
	}
}

func (s *loggingUserService) CreateUser(ctx context.Context, user *services.User) error {
	start := time.Now()
	defer func() {
		s.logger.Info("CreateUser completed",
			"duration", time.Since(start),
			"user_name", user.Name,
		)
	}()
	
	s.logger.Info("CreateUser started", "user_name", user.Name)
	return s.next.CreateUser(ctx, user)
}

func (s *loggingUserService) GetUser(ctx context.Context, id uint) (*services.User, error) {
	start := time.Now()
	defer func() {
		s.logger.Info("GetUser completed",
			"duration", time.Since(start),
			"user_id", id,
		)
	}()
	
	s.logger.Info("GetUser started", "user_id", id)
	return s.next.GetUser(ctx, id)
}

func (s *loggingUserService) ListUsers(ctx context.Context) ([]*services.User, error) {
	start := time.Now()
	defer func() {
		s.logger.Info("ListUsers completed",
			"duration", time.Since(start),
		)
	}()
	
	s.logger.Info("ListUsers started")
	return s.next.ListUsers(ctx)
}
