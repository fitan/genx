package app

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/fitan/genx/plugs/fx/example/services"
)

// @fx(type="invoke")
func InitDatabase(db *gorm.DB, logger *slog.Logger) error {
	logger.Info("initializing database")

	// 自动迁移数据库表
	err := db.AutoMigrate(&services.User{})
	if err != nil {
		logger.Error("failed to migrate database", "error", err)
		return err
	}

	logger.Info("database initialized successfully")
	return nil
}

// @fx(type="invoke")
func StartApplication(
	lifecycle fx.Lifecycle,
	userSvc services.UserService,
	logger *slog.Logger,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("application starting")

			// 创建一个测试用户
			testUser := &services.User{
				Name:  "Test User",
				Email: "test@example.com",
			}

			err := userSvc.CreateUser(ctx, testUser)
			if err != nil {
				logger.Error("failed to create test user", "error", err)
				return err
			}

			// 获取用户列表
			users, err := userSvc.ListUsers(ctx)
			if err != nil {
				logger.Error("failed to list users", "error", err)
				return err
			}

			logger.Info("users in database", "count", len(users))
			for _, user := range users {
				logger.Info("user", "id", user.ID, "name", user.Name, "email", user.Email)
			}

			logger.Info("application started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("application stopping")
			return nil
		},
	})
}
