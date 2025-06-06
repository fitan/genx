package services

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
)

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id uint) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
}

type userServiceImpl struct {
	db     *gorm.DB
	logger *slog.Logger
}

// @fx(type="provide")
func NewUserService(db *gorm.DB, logger *slog.Logger) UserService {
	return &userServiceImpl{
		db:     db,
		logger: logger,
	}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, user *User) error {
	s.logger.Info("creating user", "name", user.Name)
	return s.db.WithContext(ctx).Create(user).Error
}

func (s *userServiceImpl) GetUser(ctx context.Context, id uint) (*User, error) {
	s.logger.Info("getting user", "id", id)
	var user User
	err := s.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (s *userServiceImpl) ListUsers(ctx context.Context) ([]*User, error) {
	s.logger.Info("listing users")
	var users []*User
	err := s.db.WithContext(ctx).Find(&users).Error
	return users, err
}

// UserHandler HTTP处理器
type UserHandler struct {
	svc UserService
}

// @fx(type="provide", group="handlers")
func NewUserHandler(svc UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) CreateUser(ctx context.Context, user *User) error {
	return h.svc.CreateUser(ctx, user)
}

func (h *UserHandler) GetUser(ctx context.Context, id uint) (*User, error) {
	return h.svc.GetUser(ctx, id)
}

func (h *UserHandler) ListUsers(ctx context.Context) ([]*User, error) {
	return h.svc.ListUsers(ctx)
}
